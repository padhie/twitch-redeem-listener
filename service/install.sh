#!/bin/bash

# Pfade auflösen
# Ermittelt den absoluten Pfad zum Verzeichnis, in dem dieses Skript liegt
SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
# Das Projekt-Verzeichnis ist eine Ebene über dem service-Verzeichnis
PROJECT_DIR=$(cd "$SCRIPT_DIR/.." && pwd)
SERVICE_FILE_SRC="$SCRIPT_DIR/twitch-redeem.service"
SERVICE_FILE_DEST="/etc/systemd/system/twitch-redeem.service"

# Variablen für Ersetzung
CURRENT_USER=$(whoami)

echo "Installing Twitch Redeem Trigger Service..."
echo "User: $CURRENT_USER"
echo "Directory: $PROJECT_DIR"

# Prüfen ob Source existiert
if [ ! -f "$SERVICE_FILE_SRC" ]; then
    echo "Fehler: $SERVICE_FILE_SRC nicht gefunden!"
    exit 1
fi

# Ersetzung durchführen und in temporäre Datei schreiben
TEMP_FILE=$(mktemp)
sed "s|{USER}|$CURRENT_USER|g; s|{DIRECTORY}|$PROJECT_DIR|g" "$SERVICE_FILE_SRC" > "$TEMP_FILE"

# In Zielverzeichnis kopieren (erfordert sudo)
echo "Kopiere Service-Datei nach $SERVICE_FILE_DEST..."
if sudo cp -f "$TEMP_FILE" "$SERVICE_FILE_DEST"; then
    rm "$TEMP_FILE"
    echo "Service erfolgreich installiert (bestehende Datei wurde überschrieben)."
    
    # Systemd neu laden
    echo "Lade systemd neu..."
    sudo systemctl daemon-reload
    
    # Falls der Service schon läuft, neu starten
    if systemctl is-active --quiet twitch-redeem; then
        echo "Dienst läuft bereits, starte neu um Änderungen zu übernehmen..."
        sudo systemctl restart twitch-redeem
    fi
else
    echo "Fehler beim Kopieren der Datei. Hast du sudo-Rechte?"
    rm "$TEMP_FILE"
    exit 1
fi
