# pip install azure-cognitiveservices-vision-computervision
from azure.cognitiveservices.vision.computervision import ComputerVisionClient
from azure.cognitiveservices.vision.computervision.models import OperationStatusCodes
from msrest.authentication import CognitiveServicesCredentials
import os
import sys
import time

# Authenticate
subscription_key = "YOUR_SUBSCRIPTION_KEY"
endpoint = "YOUR_ENDPOINT"
computervision_client = ComputerVisionClient(endpoint, CognitiveServicesCredentials(subscription_key))

# Analyze an image
image_url = "https://example.com/sample.jpg"

features = ['categories', 'description', 'color']
details = ['celebrities']
description_details = ['captions']

results = computervision_client.analyze_image(image_url, visual_features=features, details=details, description_details=description_details)

if 'request_id' in results.headers:
    # Check for Text in Image
    if 'description' in results.analyze_result:
        print("Image Description:")
        for caption in results.analyze_result['description']['captions']:
            print(f"{caption['text']} (Confidence: {caption['confidence']})")

    # Display Categories
    if 'categories' in results.analyze_result:
        print("Image Categories:")
        for category in results.analyze_result['categories']:
            print(f"{category['name']} (Confidence: {category['score']})")

    # Display Colors
    if 'color' in results.analyze_result:
        print("Dominant Colors:")
        for color in results.analyze_result['color']['dominant_colors']:
            print(color)

    # Detect Celebrities
    if 'celebrities' in results.analyze_result:
        print("Celebrities:")
        for celeb in results.analyze_result['celebrities']:
            print(celeb['name'])

else:
    print("Image analysis failed. Status Code: ", results.status_code)