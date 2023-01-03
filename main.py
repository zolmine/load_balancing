# import smtplib
# from email.mime.multipart import MIMEMultipart
# from email.mime.text import MIMEText

# print("sending")
# msg = MIMEMultipart('alternative')
# msg['Subject'] = 'An example email'
# msg['From'] = 'first.last@gmail.com'
# msg['To'] = 'amine@thebay.ma'

# print("all setted")

# part1 = MIMEText("Hello!", 'plain')
# part2 = MIMEText("<h1>Hello!</h1>", 'html')

# msg.attach(part1)
# msg.attach(part2)

# # Send the message via our own SMTP server.
# s = smtplib.SMTP('mail.thebay.ma', 465)
# # s = smtplib.SMTP('mail.thebay.ma', port=465)
# print(s)
# s.send_message(msg)
# s.quit()


# import smtplib
# from email.mime.text import MIMEText

# sender = 'amine@thebay.ma'
# receivers = ['amine-es-sobhi@live.com']


# port = 993
# msg = MIMEText('This is test mail')

# msg['Subject'] = 'Test mail'
# msg['From'] = 'amine@thebay.ma'
# msg['To'] = 'amine-es-sobhi@live.com'


# with smtplib.SMTP("smtp.mailtrap.io", 2525) as server:
#     server.login("4e786c31d1e945", "e2ffb74a5968ab")
#     server.sendmail(sender, receivers, msg.as_string())
#     print("Successfully sent email")


import smtplib

sender = "amine@thebay.ma"
receiver = "elhour.yousef@gmail.com"

message = f"""\
Subject: Hi Mailtrap
To: {receiver}
From: {sender}

This is a test e-mail message."""

with smtplib.SMTP("smtp-relay.sendinblue.com", 587) as server:
    server.login("amine@thebay.ma", "gdJYBCm1P2LjswSp")
    # server.login("karama", "xsmtpsib-2063dde40d861c80ab7fed76c5071c9e83e23877417f404f48943237782165c3-MjHmPdQ5X6xRDfk0@wL1l")
    # server.
    server.sendmail(sender, receiver, message)

# ------------------
# Create a campaign\
# ------------------
# Include the Sendinblue library\
# from __future__ import print_function
# import time
# import sib_api_v3_sdk
# from sib_api_v3_sdk.rest import ApiException
# from pprint import pprint
# # Instantiate the client\
# sib_api_v3_sdk.configuration.api_key['api-key'] = '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'
# api_instance = sib_api_v3_sdk.EmailCampaignsApi()
# # Define the campaign settings\
# email_campaigns = sib_api_v3_sdk.CreateEmailCampaign(
# name= "Campaign sent via the API",
# subject= "My subject",
# sender= { "name": "From name", "email": "amine@thebay.ma"},
# type= "classic",
# # Content that will be sent\
# html_content= "Congratulations! You successfully sent this example campaign via the Sendinblue API.",
# # Select the recipients\
# recipients= {"listIds": [2, 7]},
# # Schedule the sending in one hour\
# # scheduled_at= "2022-12-01 00:00:01"
# )
# # Make the call to the client\
# try:
#     api_response = api_instance.create_email_campaign(email_campaigns)
#     pprint(api_response)
# except ApiException as e:
#     print("Exception when calling EmailCampaignsApi->create_email_campaign: %s\n" % e)