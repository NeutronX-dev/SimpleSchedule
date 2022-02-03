@echo off
cd ..
echo +---------------- SimpleSchedule ----------------+
echo ^|          NeutronX-dev/SimpleSchedule           ^|
echo ^|                                                ^|
echo +------------------- OS List --------------------+
echo ^|    android  ^|  android/arm  ^|  android/arm64   ^|
echo ^|     linux   ^|    window     ^|  android/darwin  ^|
echo ^|                                                ^|
echo +-------------- Build Information ---------------+
set /p version="| [92mVERSION[0m="
set /p os="| [92mOS[0m="
echo +----------------- Other Info -------------------+
echo Logo: SimpleSchedule\logos\noBG-1000x1000-SimpleSchedule.png
echo +---------------- Build Output ------------------+
fyne package -name "SimpleSchedule" -icon "%cd%\logos\noBG-1000x1000-SimpleSchedule.png" -appVersion "%version%" -os "%os%"
echo +------------------------------------------------+
pause