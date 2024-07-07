# Copyright (c) 2024 Charles Ozochukwu

# Permission is hereby granted, free of charge, to any person obtaining a copy
#  of this software and associated documentation files (the "Software"), to deal
#  in the Software without restriction, including without limitation the rights
#  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
#  copies of the Software, and to permit persons to whom the Software is
#  furnished to do so, subject to the following conditions:

# The above copyright notice and this permission notice shall be included in all
#  copies or substantial portions of the Software.

#  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
#  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
#  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
#  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
#  SOFTWARE
"""
This module defines the function that interacts with the g4f package.
"""
from enum import Enum
from typing import List
import g4f
from g4f.Provider import RetryProvider, GeminiPro
import g4f.image


class Message(str, Enum):
    """ Message enum with role and content attributes """
    role: str
    content: str


_PROVIDERS = [
    g4f.Provider.Liaobots,
    g4f.Provider.Feedough,
    g4f.Provider.Cnote,
    g4f.Provider.Blackbox,
    g4f.Provider.Aichatos,
    g4f.Provider.Llama,
    g4f.Provider.HuggingFace,
    g4f.Provider.Koala,
    g4f.Provider.Replicate,
    g4f.Provider.Groq,
    g4f.Provider.FreeChatgpt,
    g4f.Provider.MetaAI,
]

g4f.debug.logging = True


def _create_message_list(messages: List[Message]) -> List[dict]:
    """Create a list of dictionaries from Message enum"""
    return [{"role": message.role, "content": message.content} for message in messages]


def _handle_stream_response(response: iter) -> List[str]:
    """Handle stream response and return a list of strings"""
    return list(response)


def get_chat_completion_response(
    grok_api_key: str,
    gemini_api_key: str,
    messages: List[Message],
    model: str = None,
    stream: bool = False,
    proxy: str = None,
    timeout: int = None,
    shuffle: bool = False,
    image: bytes = None
):
    """
    Get chat completion response from g4f library

    Args:
        messages: List of Message enum
        api_key: API key for g4f library
        proxy: Proxy URL for g4f library
        stream: Whether to stream the response
        timeout: Timeout for g4f library
        model: Model to use for chat completion
        shuffle: Whether to shuffle providers


    """
    try:
        message_list = _create_message_list(messages)

        if model is not None and model not in ["gpt-3.5-turbo", "gpt-4"]:
            raise ValueError(
                "Invalid model. Only 'gpt-3.5-turbo' or 'gpt-4' are allowed.")
        if model in ["gpt-3.5-turbo", "gpt-4"]:
            response = g4f.ChatCompletion.create(
                model=model,
                messages=message_list,
                stream=stream,
                proxy=proxy,
                timeout=timeout
            )
        if image is not None:
            response = g4f.ChatCompletion.create(
                provider=GeminiPro,
                model=model,
                image=image,
                api_key=gemini_api_key,
                messages=message_list,
                stream=stream,
                proxy=proxy,
                timeout=timeout,
            )
        else:
            response = g4f.ChatCompletion.create(
                model=g4f.models.default,
                provider=RetryProvider(_PROVIDERS, shuffle=shuffle),
                api_key=grok_api_key,
                messages=message_list,
                stream=stream,
                proxy=proxy,
                timeout=timeout,
            )

        if stream:
            return _handle_stream_response(response)

        return response

    except ValueError as e:
        if model == "gpt-4":
            print("Fallback to 'gpt-3.5-turbo'")
            return get_chat_completion_response(
                messages, grok_api_key, proxy, stream, timeout, model="gpt-3.5-turbo"
            )
        raise e

    except Exception as e:
        raise e

# def handle_image_processing(
#     image: bytes,
#     messages: List[Message],
#     api_key: str = None,
#     proxy: str = None,
#     stream: bool = False,
#     timeout: int = None,
#     model: str = None
# ):
#     """
#     Handle image processing for g4f library

#     Args:
#         image: Image in bytes format
#         messages: List of Message enum
#         api_key: API key for g4f library
#         proxy: Proxy URL for g4f library
#         stream: Whether to stream the response
#         timeout: Timeout for g4f library
#         model: Model to use for chat completion


#     """
#     try:
#         message_list = _create_message_list(messages)

#         response = g4f.ChatCompletion.create(
#             model=model,
#             provider=GeminiPro,
#             image=image,
#             api_key=api_key,
#             messages=message_list,
#             stream=stream,
#             proxy=proxy,
#             timeout=timeout,
#         )

#         if stream:
#             return _handle_stream_response(response)
#         else:
#             return response

#     except Exception as e:
#         raise Exception("Error occurred while handling image processing") from e
