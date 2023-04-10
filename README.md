# container-restarter

![alt text](https://github.com/sooraj-sky/container-restarter/blob/main/assets/images/Screenshot.png?raw=true)


This application provides a user-friendly graphical interface for managing Docker containers. It uses the Fyne toolkit for Go, which allows for the creation of beautiful and responsive GUI applications. With this application, you can view a list of your Docker containers and their status, as well as start, stop, and remove containers with just a few clicks.

Moreover, the application includes a feature that monitors a specified folder for changes in its contents. This is particularly useful when developing applications that run inside Docker containers, as it allows you to automatically restart the container when changes are made to the code. The application will detect changes in the specified folder and automatically restart the Docker container associated with the project in that folder, saving you the hassle of having to manually stop and start the container every time you make changes to your code.

Overall, this application makes it easy to manage Docker containers and streamlines the development process by automating the task of restarting containers when changes are made to the code.