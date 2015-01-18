# Timeplaner for Universitetet i Agder
Dette er et program for å generere .csv-filer fra timeplansystemet til UiA. Disse filene kan så importeres direkte til Google Calendar, og synkroniseres mot din smarttelefon eller PC.

Ferdig genererte filer ligger i "timeplaner"-mappen.

### Hvordan bruker jeg timeplanen i Google Calendar?
Se [denne](https://support.google.com/calendar/answer/37118?hl=en) siden. Kort forklart oppretter du en kalender, for så å trykke "Andre kalendere" og "Importer kalender". Her velger du kalenderen du ønsker at timeplanen skal legges til i, og velger .csv-filen herfra (åpne filen på Github, trykk "Raw", og lagre filen på maskinen).

### Hvordan kan jeg generere nye versjoner av timeplanene?
Installer [Go](http://golang.org/doc/install), og kjør "go get github.com/oal/timeplan-uia", og kjør så "timeplan-uia" fra din `$GOPATH/bin`.

Alternativt kan du klone denne repositoryen, og kjøre `go run timeplan.go` fra denne mappen.

### Hvorfor trenger jeg dette?
Jeg lagde dette programmet for å spare litt tid selv med å legge inn min kalender på dataingeniørstudiet, men det var såpass lite ekstra arbeid å utvide det til andre linjer også, så hvorfor ikke. Frem til UiA offisielt får støtte for iCal/CSV, vil dette verktøyet forhåpentligvis være nyttig for noen.