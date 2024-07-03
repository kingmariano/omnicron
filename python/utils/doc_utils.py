async def process_page(page):
            tp = page.get_textpage_ocr(full=True)
            return page.get_text(textpage=tp)