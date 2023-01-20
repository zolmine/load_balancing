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


# import smtplib

# sender = "amine@thebay.ma"
# receiver = "elhour.yousef@gmail.com"

# message = f"""\
# Subject: Hi Mailtrap
# To: {receiver}
# From: {sender}

# This is a test e-mail message."""

# with smtplib.SMTP("mail.thebay.ma", 465) as server:
#     server.login("amine@thebay.ma", "gdJYBCm1P2LjswSp")
#     # server.login("karama", "xsmtpsib-2063dde40d861c80ab7fed76c5071c9e83e23877417f404f48943237782165c3-MjHmPdQ5X6xRDfk0@wL1l")
#     # server.
#     server.sendmail(sender, receiver, message)

from csv import reader
import smtplib
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart
from email.mime.base import MIMEBase
from email import encoders
from email.mime.image import MIMEImage
from time import sleep

file = "mails.csv"
f = open("mails.csv","r")
data = reader(f)
# print(data)

email_user = 'amine@thebay.ma'
email_password = 'xsmtpsib-2063dde40d861c80ab7fed76c5071c9e83e23877417f404f48943237782165c3-MJZvbwaTk6hIq7U8'
server = smtplib.SMTP('smtp-relay.sendinblue.com',587)
server.starttls()
server.login(email_user,email_password)

for item in data:
  # print(item)
  email_send = "smailb@thebay.ma"
  subject = f"Let's talk about {item[2]}'s future !"
  msg = MIMEMultipart()
  msg['From'] = "smailb@thebay.ma"
  msg['To'] = item[1]
  msg['Subject'] = subject
  h = item[0]
  c = item[2]
  # t = str("trainer")
  html = """\
  <div style="font-family: 'Poppins', sans-serif; margin-bottom: 40px">
    <p>Dear {first_name},</p>
    <p>I hope this email finds you well, I am Smail, founder, and CEO of The Tech Bay, I came across {startup_name} online, and what you guys are doing picked my interest!</p>
    <p>I know firsthand how challenging it can be for a startup to get off the ground, especially when it comes to developing and marketing innovative products. My team and I have years of experience working with startups like yours turning their ideas into scale-ups in complex environments..</p>
    <p>{first_name}, I'm sure we can provide you with valuable information and support for your business development.</p>
    <p>Let's e-meet!  <a href="https://calendly.com/thetechbay/smail" target="_blank">[here]</a></p>
    <p>Cheers.</p>
    <p>Smail E</p>
    </div>

  <div style="display: flex; flex-direction: column; font-family: 'Poppins', sans-serif;">
  <p style="margin: 4px; "><strong>Moulay-Smail EL BOUKFAOUI</strong></p>
  <p style="margin: 4px; margin-bottom: 15px; ">Founder, Strategy Guy</p>

  <img style="width: 150px" src="https://d1muf25xaso8hp.cloudfront.net/https%3A%2F%2Fs3.amazonaws.com%2Fappforest_uf%2Ff1669368894661x937951137529545300%2FFINAL_PROP1_GLOBAL_Signature.png?w=256&h=58&auto=compress&dpr=1&fit=max" alt="Company logo">
  <p style="margin: 4px; "><a href="tel:+212700175518" target="_blank">+212700175518</a> | <a href="tel: +17869336054"> +17869336054</a></p>
  <p style="margin: 4px; "><a href="mailto:smailb@thebay.ma" target="_blank">smailb@thebay.ma</a></p>
  <p style="margin: 4px; "><a href="https://www.thebay.ma/" target="_blank">www.thebay.ma</a></p>
</div>
  """.format(first_name=h , startup_name=c) 
      
  msg.attach(MIMEText(html,'html'))

  server.sendmail(email_send,item[1],msg.as_string())
  # sleep()
  
server.quit()

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