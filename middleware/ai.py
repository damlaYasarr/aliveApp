#!/usr/bin/python
import pathlib
import textwrap
import google.generativeai as genai
from dotenv import load_dotenv
import os
import sys
import json
from IPython.display import display, Markdown

load_dotenv()

GOOGLE_API_KEY = os.getenv("GEMINI_API_KEY")
genai.configure(api_key=GOOGLE_API_KEY)

model = genai.GenerativeModel('gemini-1.5-flash')

# Function to format text into Markdown
def to_markdown(text):
    text = text.replace('•', '  *')
    return Markdown(textwrap.indent(text, '> ', predicate=lambda _: True))

# Function to query the AI
def ask_anything(text):
    response = model.generate_content(text)
    return response.text

# Main block to take input and output the AI feedback in JSON format
if __name__ == "__main__":
   
    feedback = sys.argv[1]
    
    # Generate AI feedback
    ai_response = ask_anything(feedback)
    
    # Return the response as JSON
    print(json.dumps({"response": ai_response}))