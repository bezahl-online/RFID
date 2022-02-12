// Basic MFRC522 RFID Reader Code by cooper @ my.makesmart.net
// Released under Creative Commons - CC by cooper@my.makesmart.net

#include <SPI.h>
#include <MFRC522.h>

#define SS_PIN          15         // Pin: D8
#define RST_PIN         0          // Pin: D3

double chipID;


// MFRC522-Instanz erstellen
MFRC522 mfrc522(SS_PIN, RST_PIN);

void setup() {
  Serial.begin(9600);

  // SPI-Bus initialisieren
  SPI.begin();

  // MFRC522 initialisieren
  mfrc522.PCD_Init();

  //Kurze Pause nach dem Initialisieren   
  delay(10);
  Serial.println("");

  // Details vom MFRC522 RFID READER / WRITER ausgeben
  mfrc522.PCD_DumpVersionToSerial();  

  
  Serial.println("RFID-Chip auflegen, um UID anzuzeigen...");
}

void loop() {
  
  // Sobald ein Chip aufgelegt wird startet diese Abfrage
  if (mfrc522.PICC_IsNewCardPresent()){
  
    //CardID resetten
    chipID = 0;
    
    mfrc522.PICC_ReadCardSerial();

 char str[32] = "";
   UID_Sender(mfrc522.uid.uidByte, 7, str); //Insert (byte array, length, char array for output)
   Serial.println(str);
   
   /* Hier wird die ID des Chips in die Variable chipID geladen
    for (byte i = 0; i < mfrc522.uid.size; i++){
      
    }

    //... und anschließend ausgegeben
    Serial.println(chipID);
*/
    // Danach 5 Sekunden pausieren, um mehrfaches lesen /ausführen zu verhindern
    delay(500);
      
  }
}

void UID_Sender(byte array[], unsigned int len, char buffer[])
{
   for (unsigned int i = 0; i < len; i++)
   {
      byte nib1 = (array[i] >> 4) & 0x0F;
      byte nib2 = (array[i] >> 0) & 0x0F;
      buffer[i*2+0] = nib1  < 0xA ? '0' + nib1  : 'A' + nib1  - 0xA;
      buffer[i*2+1] = nib2  < 0xA ? '0' + nib2  : 'A' + nib2  - 0xA;
   }
   buffer[len*2] = '\0';
}