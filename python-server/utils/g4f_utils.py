import g4f
from g4f.Provider import RetryProvider, GeminiPro
from typing import List
from enum import Enum
class Message(str, Enum):
    role: str
    content: str
    
_providers = [
    g4f.Provider.DuckDuckGo,
    g4f.Provider.Blackbox,
    g4f.Provider.Cnote,
    g4f.Provider.HuggingFace,
    g4f.Provider.Ecosia,
    g4f.Provider.Liaobots,
    g4f.Provider.Feedough,
    g4f.Provider.Koala,
    g4f.Provider.Llama,
    g4f.Provider.Replicate,
    g4f.Provider.Groq,
    g4f.Provider.MetaAI,
   g4f.Provider.FreeChatgpt,
]
g4f.debug.logging = True

def handle_stream_response(response: str):
    streamResponse = []
    streamResponse.append(response)
    return streamResponse


def get_chat_completion_response(messages: List[Message], api_key: str = None, proxy: str = None, stream: bool = False, timeout: int = None, model: str = None, shuffle: bool = False):
    response = ""
    try:
        # Convert messages to a list of dictionaries
        message_list = [{"role": message.role, "content": message.content} for message in messages]

        if model is not None and model not in ["gpt-3.5-turbo", "gpt-4"]:
            raise ValueError("Invalid model. Only 'gpt-3.5-turbo' or 'gpt-4' are allowed.")
        
        if model in ["gpt-3.5-turbo", "gpt-4"]:
            response = g4f.ChatCompletion.create(
                model=model,
                messages=message_list,
                stream=stream,
                proxy=proxy,
                timeout=timeout
            )
        else:
            response = g4f.ChatCompletion.create(
                model=g4f.models.default,
                provider=RetryProvider(_providers, shuffle=shuffle),
                api_key=api_key,
                messages=message_list,
                stream=stream, 
                proxy=proxy,
                timeout=timeout,
            
                
            )

        if stream:
            # Handle streaming response
            streamResponse = []
            for message in response:
                streamResponse.append(message)
            return streamResponse    
        else:
            # Handle non-streaming response
            # Process the response as needed
            pass

    except Exception as e:
        # If an error occurs, handle it here
        if model == "gpt-4":
            print("Fallback to 'gpt-3.5-turbo'")
            # Retry with a different model
            return get_chat_completion_response(messages, api_key, proxy, stream, timeout, model="gpt-3.5-turbo")
        raise Exception(e)

    return response


def handle_image_processing(image: bytes, messages: List[Message], api_key: str = None, proxy: str = None, stream: bool = False, timeout: int = None, model: str = None):
    response = ""
    try:
        message_list = [{"role": message.role, "content": message.content} for message in messages]
     
        response = g4f.ChatCompletion.create(
            model = model,
            provider= GeminiPro,
            image = image,
            api_key=api_key,
            messages=message_list,
            stream=stream, 
            proxy=proxy,
            timeout=timeout,
        ) 
        
        if stream:
            # Handle streaming response
            streamResponse = []
            for message in response:
                streamResponse.append(message)
            return streamResponse    
        else:
            # Handle non-streaming response
            # Process the response as needed
            pass   
        
    except Exception as e:

        raise Exception(e)
    return response    