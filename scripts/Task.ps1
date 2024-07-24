try{
    $arguments = get-content "C:\Admin\ConfigData\data.txt" -Raw -ErrorAction Stop
    Remove-Item "C:\Admin\ConfigData\data.txt" 
    $x = $arguments -split "-"
    $x = $x[1] -split" "
    $x = $x[1]+".cfg"
    if(-not (Test-Path ("C:\inetpub\oc_configurator\configs\"+$x) -PathType Leaf) ){
           cd C:\inetpub\oc_configurator\configs
           Start-Process .\OctopusConfigurator.exe -ArgumentList $arguments
    }
    } catch{
        Write-Host "No Data File avaiable"
    } 

# Registerd as a task on the server running as a service account