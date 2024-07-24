#Declarejob that enables and disables the gmsa after 10 minutes
# Define the job trigger to run every minute
$T= New-JobTrigger -Once -At (Get-Date).AddMinutes(1) -RepetitionInterval (New-TimeSpan -Minutes 1) -RepetitionDuration ([TimeSpan]::MaxValue)
$O = New-ScheduledJobOption -RunElevated -MultipleInstancePolicy "IgnoreNew" -RequireNetwork -StartIfOnBattery -DoNotAllowDemandStart
      
Register-ScheduledJob -Name "ReadNewConfig" -Trigger $T -ScheduledJobOption $O -ScriptBlock {
    try{
        $arguments = get-content "C:\Admin\ConfigData\data.txt" -Raw -ErrorAction Stop
        $path = "C:\inetpub\oc_configurator\configs\OctopusConfigurator.exe"
        Start-Process -FilePath $path -ArgumentList $arguments -Verb RunAs
        Remove-Item "C:\Admin\ConfigData\data.txt"
        } catch{}
    }     
      

#Set so GMSA runs the exe
$P = New-ScheduledTaskPrincipal -UserID "NT AUTHORITY\SYSTEM" -LogonType ServiceAccount -RunLevel Highest;
$psJobsPathInScheduler = "\Microsoft\Windows\PowerShell\ScheduledJobs";
Set-ScheduledTask -TaskPath $psJobsPathInScheduler -TaskName "ReadNewConfig" -Principal $P
Write-Host "Task created under Microsoft\Windows\PowerShell\ScheduledJobs" -ForegroundColor Green
      
