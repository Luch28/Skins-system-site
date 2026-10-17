# Skins-system-site
A simple website for your own skin system

He works with databases to store player data. I have attached the SQL files

## To configure
In the launcher, to configure requests on the website:
1. open the file ```LauncherFrame.java```
2. find:
 ```java
 String result = Util.excutePost("http://87.255.8.205/auth", parameters);
```
3. enter your URL

To configure on website:
1. Create a database using the files from the sql folder
2. install the github.com/go-sql-driver/mysql package
3. go to the project folder and enter ```go run .```

The native files are already in the site's folders
> [!TIP]
> If you are using Linux, don't forget to create an auto-start service
