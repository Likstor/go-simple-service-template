Get-Content ..\configs\.env | ForEach-Object {
    $name, $value = [regex]::split($_, '[ ]*=[ ]*')
    if (!([string]::IsNullOrWhiteSpace($name) -OR $name.Contains('#'))) {
        Set-Content env:$name $value
    }
}

docker compose -f .\docker-compose.yaml up -d --build