@echo on
echo Build and Deploy
go build
echo Deploy exe and ini for festa junina web site
copy festajuninaweb.exe Z:\I_Daniel\OneDrive\I_Projects\Year2017\golang\runtime\festajunina\fjwebsite
copy festajunina.ini Z:\I_Daniel\OneDrive\I_Projects\Year2017\golang\runtime\festajunina\fjwebsite
echo Deployment is done.
