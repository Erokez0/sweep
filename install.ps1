# PowerShell installer for Windows

$GITHUB_USERNAME = "erokez0"
$APP_NAME = "sweep"

Write-Host "Installing $APP_NAME" -ForegroundColor Green

# Windows-specific settings
$OS = "windows"
$EXT = ".exe"

# Detect architecture
$arch = $env:PROCESSOR_ARCHITECTURE
switch ($arch) {
    { $_ -in "AMD64", "x64" } { $ARCH = "amd64" }
    "ARM64" { $ARCH = "arm64" }
    default {
        Write-Host "Unsupported architecture: $arch" -ForegroundColor Red
        exit 1
    }
}

Write-Host "Detected Architecture: $ARCH" -ForegroundColor Cyan

# Set directories
$BIN_DIR = "$env:APPDATA\Programs\$APP_NAME"
$CONFIG_DIR = "$env:APPDATA\$APP_NAME"
$BINARY_NAME = "$APP_NAME$EXT"

Write-Host "Binary directory: $BIN_DIR"
Write-Host "Config directory: $CONFIG_DIR"

# Create directories
New-Item -ItemType Directory -Force -Path $BIN_DIR, $CONFIG_DIR | Out-Null

# Get version
$REPO_URL = "https://github.com/$GITHUB_USERNAME/$APP_NAME"
$VERSION = if ($env:VERSION) { $env:VERSION } else { "latest" }

if ($VERSION -eq "latest") {
    try {
        $releaseInfo = Invoke-RestMethod -Uri "https://api.github.com/repos/$GITHUB_USERNAME/$APP_NAME/releases/latest"
        $VERSION = $releaseInfo.tag_name
    } catch {
        Write-Host "Warning: Could not fetch latest version, using v1.0.0" -ForegroundColor Yellow
        $VERSION = "v1.0.0"
    }
}

# Download binary
$BINARY_URL = "$REPO_URL/releases/download/$VERSION/${APP_NAME}-${OS}-${ARCH}${EXT}"
$TEMP_BINARY = "$env:TEMP\$BINARY_NAME"

Write-Host "Downloading binary..." -ForegroundColor Cyan

try {
    Invoke-WebRequest -Uri $BINARY_URL -OutFile $TEMP_BINARY -ErrorAction Stop
} catch {
    Write-Host "Failed to download binary: $_" -ForegroundColor Red
    exit 1
}

# Verify download
if (-not (Test-Path $TEMP_BINARY)) {
    Write-Host "Download failed - file not found" -ForegroundColor Red
    exit 1
}

$fileSize = (Get-Item $TEMP_BINARY).Length
if ($fileSize -lt 1000) {
    Write-Host "Download failed - file too small ($fileSize bytes)" -ForegroundColor Red
    exit 1
}

Write-Host "Binary downloaded ($([math]::Round($fileSize/1KB, 2)) KB)" -ForegroundColor Green

# Install binary
$INSTALL_PATH = "$BIN_DIR\$BINARY_NAME"
Copy-Item -Path $TEMP_BINARY -Destination $INSTALL_PATH -Force

# Clean up
Remove-Item -Path $TEMP_BINARY -Force -ErrorAction SilentlyContinue

# Download config files if they don't exist
$CONFIG_BASE_URL = "$REPO_URL/raw/main/"
$CONFIG_FILES = @("config.schema.json", "config.default.json")

foreach ($file in $CONFIG_FILES) {
    $CONFIG_PATH = "$CONFIG_DIR\$file"
    
    if (-not (Test-Path $CONFIG_PATH)) {
        try {
            Invoke-WebRequest -Uri "$CONFIG_BASE_URL/$file" -OutFile $CONFIG_PATH -ErrorAction Stop
        } catch {
            Write-Host "Could not download $file, creating empty file" -ForegroundColor Yellow
            New-Item -ItemType File -Path $CONFIG_PATH -Force | Out-Null
        }
    }
}

# Create user config from default if it doesn't exist
$DEFAULT_CONFIG = "$CONFIG_DIR\config.default.json"
$USER_CONFIG = "$CONFIG_DIR\config.json"

if ((Test-Path $DEFAULT_CONFIG) -and (-not (Test-Path $USER_CONFIG))) {
    Copy-Item -Path $DEFAULT_CONFIG -Destination $USER_CONFIG -Force
}

# Add to PATH
Write-Host "Setting up PATH..." -ForegroundColor Green

$currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($currentPath -split ';' -notcontains $BIN_DIR) {
    # Check if it's already in PATH as a subdirectory
    $alreadyInPath = $false
    foreach ($pathEntry in ($currentPath -split ';')) {
        if ($pathEntry -and (Test-Path $pathEntry) -and $BIN_DIR.StartsWith($pathEntry)) {
            $alreadyInPath = $true
            break
        }
    }
    
    if (-not $alreadyInPath) {
        [Environment]::SetEnvironmentVariable("PATH", "$currentPath;$BIN_DIR", "User")
        Write-Host "Added $BIN_DIR to user PATH" -ForegroundColor Green
        Write-Host "You may need to restart your terminal for changes to take effect" -ForegroundColor Yellow
    } else {
        Write-Host "$BIN_DIR is already accessible via PATH" -ForegroundColor Cyan
    }
} else {
    Write-Host "$BIN_DIR is already in PATH" -ForegroundColor Cyan
}

Write-Host "`nInstallation Complete!" -ForegroundColor Green
Write-Host "`nTo run $APP_NAME, open a new terminal and type:" -ForegroundColor Cyan
Write-Host "  $APP_NAME" -ForegroundColor White