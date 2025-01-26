## come giocatore, voglio poter giocare una partita contro dei bot di diversa difficoltà

- __Titolo:__ Verificare che la partita giocatore-bot venga creata correttamente
- __Obiettivo:__ Verificare che il primo giocatore sia l'utente e il secondo sia il bot scelto e che la partita venga riconosciuta come partita con il bot
- __Scenario:__ L'utente crea una partita con il bot difficile
- __Passaggi:__
   1. creare una partita utente-bot
- __Risultato atteso:__ la partita viene creata con le caratteristiche sopra elencate

- __Titolo:__ Calcolo corretto delle mosse per il bot difficile
- __Obiettivo:__ Verificare il bot difficile giochi le migliori mosse secondo l'api utilizzata
- __Scenario:__ Il bot calcola la propria mossa nella configurazione iniziale della board
- __Passaggi:__
   1. creare una partita con il bot difficile come giocatore
   2. fare calcolare la propria mossa al bot
- __Risultato atteso:__ il bot gioca la mossa uguale alla migliore mossa secondo l'api bgweb

## come creatore di tornei, voglio essere in grado di creare, amministrare tornei

- __Titolo:__ Verificare che il torneo venga creato correttamente
- __Obiettivo:__ Verificare un utente possa creare un torneo e questo sia creato correttamente
- __Scenario:__ Un utente crea un torneo
- __Passaggi:__
   1. navigare all'interfaccia di creazione tornei
   2. inserire il nome del torneo che si vuole creare
   2. creare un torneo utilizzando l'apposito bottone
- __Risultato atteso:__ un nuovo torneo deve essere creato con i seguenti campi: l'utente che lo ha creato come owner, come unico utente l'owner e come stato waiting.

- __Titolo:__ Verificare che un nuovo utente possa unirsi ad un torneo
- __Obiettivo:__ Verificare che un utente possa unirsi ad un torneo non ancora iniziato
- __Scenario:__ Un utente vuole unirsi ad un torneo
- __Passaggi:__
   1. l'utente owner crea un torneo
   2. un altro utente sceglie il torneo creato in precedenza dalla lista di tornei non ancora inizziati e si aggiunge ad esso.
- __Risultato atteso:__ L'utente deve essere aggiunto al torneo, il torneo adesso avrà due giocatori.

## come giocatore, voglio poter giocare una partita in locale

- __Titolo:__ Giocare una partita in locale
- __Obiettivo:__ Verificare che si possa creare e giocare una partita in locale
- __Scenario:__ Un utente crea una partita in locale, si effettuano due turni e si chiude la partita abbandonando
- __Passaggi:__
   1. un utente crea una partita in locale
   2. effettua il turno 1
   3. effettua il turno 2
   3. chiude la partita abbandonando
- __Risultato atteso:__ la partita viene creata con lo stesso giocatore come giocatore 1 e 2, i due turni si possono giocare dallo stesso giocatore, dopo l'abbandono la partita viene chiusa e lo stato viene messo a winp2

## come giocatore, voglio poter giocare una partita online

- __Titolo:__ Verificare che il matchmaking funzioni
- __Obiettivo:__ Verificare che tramite il matchmaking si possa trovare un avversario per le partite online
- __Scenario:__ Come giocatore voglio trovare un avversario tramite il matchmaking per giocare online
- __Passaggi:__
   1. un primo utente si mette in coda per giocare con elo 1000
   2. un secondo giocatore si mette in coda per giocare con elo 1000
- __Risultato atteso:__ si deve creare una partita online con i due giocatori

- __Titolo:__ Validazione delle mosse (manca)
- __Obiettivo:__ validare le mosse in modo da poter giocare solo le mosse permesse
- __Scenario:__ un utente si trova in partita e manda le mosse che vuole effettuare al server
- __Passaggi:__
   1. manda delle mosse errate al server
- __Risultato atteso:__ il server ritorna un messaggio di errore 400, le mosse non vengono effettuate

## come giocatore, voglio avere un interfaccia front-end per navigare all'interno del sito

- __Titolo:__ Navigabilità del sito
- __Obiettivo:__ Controllare che tutte le pagine siano renderizzate correttamente e che siano raggiungibili
- __Scenario:__ Un utente accede alla home del sito tramite un browser
- __Passaggi:__
   1. renderizzare tutte le view e i componenti
- __Risultato atteso:__ tutte le views e i componenti devono essere renderizzati correttamente

## come giocatore, voglio poter scegliere vari temi grafici e ambientazioni

- __Titolo:__ Verificare il cambio dei temi grafici
- __Obiettivo:__ Verificare il cambio tema sia rispettato in tutte le pagine
- __Scenario:__ Un utente vuole cambiare il tema del sito cambiando l'apposita impostazione
- __Passaggi:__
   1. navigare nelle impostazioni
   2. scegliere un tema proposto e applicarlo
- __Risultato atteso:__ il nuovo tema si deve applicare a tutte le pagine

## come utente, vorrei poter usare il sito anche da cellulare

- __Titolo:__ Testare l'usabilità del sito da mobile con schermo orizzontale
- __Obiettivo:__ testare che il sito si possa usare anche da dispositivi mobili
- __Scenario:__ un utente si connette al sito tramite un dispositivo mobile
- __Risultato atteso:__ tutti gli elementi nelle pagine devono essere raggiungibili, le pagine più utilizzate devono essere visualizzate senza scrolling

## come studente, voglio poter accedere a risorse di addestramento

- __Titolo:__ Modale regole
- __Obiettivo:__ Verificare la completezza e il funzionamento delle risorse di addestramento
- __Scenario:__ un utente attraverso il sito vuole imparare le regole del gioco utilizzando le risore apposite
- __Passaggi:__
   1. Collegarsi alla home del sito
   2. accedere alle risorse cliccando il bottone rules
   3. interagire con menu delle regole
- __Risultato atteso:__ la spiegazione delle regole deve essere completa e gli elementi devono essere renderizzati correttamente

## come giocatore, socievole voglio poter comunicare con i miei avversari in chat

- __Titolo:__ Invio di messaggi
- __Obiettivo:__ testare il funzionamento dell'invio di messaggi
- __Scenario:__ un utente mentre è in partita invia un messaggio all'avversario tramite l'apposita chat
- __Passaggi:__
   1. entrare in partita
   2. inviare un messaggio
- __Risultato atteso:__ il messaggio deve essere inviato al destinatario

## come influencer, voglio essere in grado di condividere i miei progressi sui social media

- __Titolo:__ Condivisione profilo su Telegram
- __Obiettivo:__ testare la corretta creazione di un messaggio su telegram che condivida il profilo
- __Scenario:__ un utente vuole condividere i suoi progressi mandando un messaggio via telegram
- __Passaggi:__
   1. accedere alla pagina statistics del sito
   2. selezionare il metodo telegram per la condivisione
   3. mandare il messaggio
- __Risultato atteso:__ il messaggio deve creato correttamente e inviato

## come utente, voglio ricevere le notifiche dal sito

- __Titolo:__ WebSockets send
- __Obiettivo:__ assicurarsi che gli utenti connessi ricevano le notifiche attraverso i websoket
- __Scenario:__ il backend manda una notifica ad un utente connesso tramite WebSockets
- __Passaggi:__
   1. mandare un messaggio di prova
- __Risultato atteso:__ il messaggio viene mandato e ricevuto correttamente

## come giocatore, voglio ricevere dei badge digitali come ricompense

- __Titolo:__ Generazione badge
- __Obiettivo:__ testare la creazione di generazione del badge prima vittoria
- __Scenario:__ un giocatore vince la prima partita
- __Passaggi:__
   1. creare una partita
   2. vincere la partita
- __Risultato atteso:__ il giocatore deve avere il badge _first victory_

## come studente, voglio essere in grado di analizzare le partite giocate

- __Titolo:__ Dati replay
- __Obiettivo:__ assicurarsi il corretto salvataggio/trasmissione dei dati di replay
- __Scenario:__ un utente dopo aver finito una partita vuole analizzare la partita fatta riguardando le mosse
- __Passaggi:__
   1. creare una partita
   2. giocare due mosse
   3. abbandonare la partita
   4. richiedere il replay delle due mosse
- __Risultato atteso:__ il dati del replay devono essere fedeli con quelli giocati in partita

## come utente, voglio tener traccia dei miei progressi nel gioco

- __Titolo:__ Dati progressi
- __Obiettivo:__ testare la generazione dei dati relativi ai progressi dei giocatori
- __Scenario:__ un utente dopo varie partite richiede le statistiche del suo profilo
- __Passaggi:__
   1. simulare una partita online
   2. simulare una partita contro un bot
   3. richiedere i dati delle statistiche
- __Risultato atteso:__ i dati delle statistiche devono essere fedeli a quelli reali

## come utente, voglio poter avere un profilo per tenere traccia delle partite giocate e delle statistiche

- __Titolo:__ Creazion del profilo
- __Obiettivo:__ testare la registrazione di un utente
- __Scenario:__ un utente vuole creare un profilo utente registrandosi al sito
- __Passaggi:__
   1. navigare alla home del sito
   2. Premere _Sing up now_
   3. Inserire tutti i campi correttamente
   4. Premere _Register_
   5. effettuare il login
- __Risultato atteso:__ l'utente deve essere creato correttamente e il login deve andare a buon fine.

