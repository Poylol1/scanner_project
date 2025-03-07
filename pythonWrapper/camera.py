from picamera2 import Picamera2 
from time import sleep
import uuid

camera = Picamera2()

config = camera.create_still_configuration(main={"size": (1920,1080)})
camera.configure(config)
camera.start_preview()
camera.start()
sleep(1)
name = str(uuid.uuid4())
imagePath = f"/home/MesaPi/Projects/scanner_project/images/${name}.jpg"
camera.capture_file(imagePath)
print(imagePath)
camera.stop_preview()
