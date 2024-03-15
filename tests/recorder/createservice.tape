Output createservice.gif
Set WindowBar Colorful
Set Height 1080
Set Width 1600
Set FontSize 20

Hide
Type "alias envme='/home/gemboran/projects/envme/bin/envme'"
Enter
Type clear
Enter
Show

Type "envme create service --interactive"
Sleep 3s
Backspace 12
Type "i"
Sleep 3s
Enter
Sleep 3s

Type "api"
Sleep 1s
Tab
Sleep 1s
Type "haqiramadhani/apicommerce"
Sleep 1s
Tab
Type "PORT=8080"
Ctrl+J
Type "DATABASE_HOST=db"
Ctrl+J
Type "DATABASE_PORT=5432"
Ctrl+J
Type "DATABASE_USER=postgres"
Ctrl+J
Type "DATABASE_PASS="
Ctrl+J
Type "DATABASE_NAME=postgres"
Tab
Sleep 1s
Type "8080:api.envme.bid"
Sleep 1s
Tab
Sleep 1s
Left
Sleep 5s
Enter
Sleep 5s
