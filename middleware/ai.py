#!/usr/bin/python
import pathlib
import textwrap

import google.generativeai as genai

from dotenv import load_dotenv
import os
import google.generativeai as genai

from IPython.display import display
from IPython.display import Markdown

load_dotenv()

GOOGLE_API_KEY = os.getenv("GEMINI_API_KEY")



genai.configure(api_key=GOOGLE_API_KEY)

model = genai.GenerativeModel('gemini-1.5-flash')
#  api key
def to_markdown(text):
  text = text.replace('•', '  *')
  return Markdown(textwrap.indent(text, '> ', predicate=lambda _: True))
def askanything(text):
      response= model.generate_content(text)
      return response.text


print(askanything("how are you doing"))


