# Retrospettiva
Team 1 - VezGammon

Valutazioni: 😃 😐 😡

## Sprint 1

| Cards            | Diego PO | Lorenzo SM | Samuele D | Fabio D | Emanuele D | Omar D |
| ---------------- | -------- | ---------- | --------- | ------- | ---------- | ------ |
| Scrum Master     | 😡       | 😡         | 😡        | 😡      | 😐         | 😐     |
| Product Owner    | 😡       | 😐         | 😐        | 😐      | 😐         | 😐     |
| Development Team | 😐       | 😡         | 😡        | 😐      | 😐         | 😐     |
| Product Backlog  | 😃       | 😃         | 😐        | 😃      | 😃         | 😃     |
| Sprint Planning  | 😡       | 😡         | 😐        | 😡      | 😡         | 😡     |
| Sprint Goal      | 😡       | 😡         | 😐        | 😐      | 😐         | 😐     |
| Self Menagement  | 😃       | 😐         | 😃        | 😃      | 😃         | 😐     |

### Commenti personali
#### Product Owner - Diego Barbieri
Sotto-stima del carico di lavoro, specialmente riguardante la tecnologia per il back-end, che porterà ad un elevato debito tecnico del secondo sprint.
Poco dialogo e male suddivisione dei compiti fra i developers. User story ben definite, task troppo vaghe. Suddivisione impari del carico di lavoro: alcuni membri hanno lavorato più di altri.
Veramente ottimo la gestione dell'ambiente di lavoro: server di deploy funzionante con un ottimo CI/CD.

#### Scrum Master - Lorenzo Peronese
Ottima gestione dell'ambiente di lavoro self-hostato, ma user stories troppo complesse e impegnative, soprattutto relative al back-end.
Cattiva organizzazione della suddivisione delle task, anche a causa mia che ho dato scadenze poco precise
e ho seguito troppo sommariamente lo sviluppo, senza aver capito a pieno lo scopo del mio ruolo.

#### Developer - Samuele Musiani
Sotto-stima del carico di lavoro, in particolare troppo per il backend e relativamente poco per il front-end.
Il team per il front si è trovato bloccato da mancanza di features nel back.

#### Developer - Fabio Murer
Sotto-stima del carico di lavoro, sprint molto impegnativo a causa di user story troppo generali richiedenti molto lavoro.
Continui problemi nel database mi hanno portato ad accumulare un ritardo nello sviluppo di features nella parte back-end.

#### Developer - Emanuele Argonni
Carico di lavoro sottostimato, abbiamo un elevato debito per il prossimo sprint. Era da pensare meglio la suddivisione delle user story considerati i 3 sprint. PO e SM dovrebbero gestire meglio il team durante l'intero periodo dello sprint evitando di assegnare attività gli ultimi giorni prima della fine dello sprint.

#### Developer - Omar Ayache
Oltre ai commenti e riflessioni dei colleghi che condivido, vorrei che al prossimo sprint si prendesse in considerazione una specie di misura (i.e. i giorni di lavoro o le ore) una misura quantitativa che associ task/tempo di lavoro. Questo dovrebbe comprendere anche la "tassa" della nuova tecnologia utilizzata, per esempio il capire come funziona la gestione del db con pSQL (oppure i possibili problemi che possono nascere).
I vantaggi dovrebbero essere che in questo modo riusciamo a valutare meglio sia il carico di lavoro ma anche il peso (punti) da assegnare ad ogni US.

### Prolematiche comuni e proposte per risolverle

Durante questo sprint, il carico di lavoro front-end è stato significativamente inferiore rispetto a quello back-end,
causando tempi di inattività per gli sviluppatori front-end in attesa del completamento delle componenti server.
Per ottimizzare gli sprint futuri, sarà fondamentale:

- Bilanciare meglio la distribuzione delle attività tra front-end e back-end
- Garantire un flusso di sviluppo continuo e efficiente

Il sovraccarico di user stories assegnate a questo sprint ha generato un debito tecnico considerevole che dovrà essere gestito
nello sprint successivo. Di conseguenza, le attività inizialmente pianificate dovranno essere riprogrammate.
Product Owner e Scrum Master si incontreranno lunedì per valutare strategie che evitino l'aggiunta di uno sprint
supplementare rispetto alla pianificazione iniziale.
Qualora l'estensione risultasse inevitabile, si procederà con una nuova distribuzione delle user stories.


## Sprint 2

Valutazioni: 😃 😐 😡

| Cards            | Diego PO | Lorenzo SM | Samuele D | Fabio D | Emanuele D | Omar D |
| ---------------- | -------- | ---------- | --------- | ------- | ---------- | ------ |
| Scrum Master     | 😃       | 😃         | 😃        | 😃      |  😃        | 😃     |
| Product Owner    | 😃       | 😃         | 😃        | 😐      |  😃        | 😃     |
| Development Team | 😐       | 😃         | 😐        | 😐      |  😃        | 😐     |
| Product Backlog  | 😃       | 😃         | 😃        | 😐      |  😐        | 😐     |
| Sprint Planning  | 😐       | 😡         | 😐        | 😃      |  😐        | 😐     |
| Sprint Goal      | 😐       | 😐         | 😐        | 😐      |  😐        | 😐     |
| Self Menagement  | 😃       | 😃         | 😃        | 😃      |  😃        | 😐     |

### Commenti personali
#### Product Owner - Diego Barbieri
Nel complesso, penso che il team abbia lavorato al meglio delle proprie capacità. Sono stati parzialmente risolti i problemi di comunicazione
e di suddivisione del carico di lavoro. Penso che sia necessaria maggiore attenzione agli standard da adottare per le API tra front e back-end.
Questo puo' essere risolto con una maggiore attenzione agli schemi di progettazione e alla documentazione.
Sono stato negativamente colpito dal fatto che durante l'ultima sprint review, alcuni membri del team non avessero ben chiaro il significato di DOD
e della sua importanza. Lo stesso e' successo quando un developer ha mergiato direttamente sul main parte dello sviluppo ancora non pienamente completato.

#### Scrum Master - Lorenzo Peronese
Il team ha lavorato di più e meglio, con maggior comunicazione e aiuto reciproco.
I daily scrum sono stati strutturati meglio e sono passati da una frequenza irregolare a una cadenza fissa ogni due giorni; sono stati un momento utile per
tutti di confronto e discussione.
Unica nota negativa è stata ancora una volta una cattiva pianificazione dello sprint: 3 user stories su 6 complessive non sono state
completate a causa di una stima sbagliata a inizio sprint e molti problemi imprevisti emersi durante queste due settimane; è stato richiesto molto impegno ai developers
ma sono soddisfatto del modo in cui queste difficoltà sono state affrontate e superate.

#### Developer - Samuele Musiani
Problemi nel back-end, soprattutto riguardo le api per backgammon. Sotto-stima dei tempi per completare certe US. Buona la scelta di spostare
un developer del frontend nel backend.

#### Developer - Emanuele Argonni
Purtroppo anche per questo sprint non siamo riusciti a portare a termine tutte la task,
dato un incoveniente con le API del backend che non ritornavano tutte le mosse possibili
Nella seconda settimana dello sprint grazie al pair-programming tra backend e frontend siamo riusciti a risolvere quasi tutti i problemi riscontrati,
nonostante il tempo perso per il debugging.
Buon lavoro dello scrum master che ha cercato di motivare i developer ed ha organizzato i daily scrum per permettere ai membri del team di confrontarsi

#### Developer - Omar Ayache
Sottostima del lavoro e problemi riguardo gestione del bot per le mosse, ottimo l'utilizzo della strategia pair-programming per i prossimi sprint.
Miglior presa di coscienza generale del carico di lavoro.

### Prolematiche comuni e proposte per risolverle
Dall'analisi delle opinioni raccolte emerge un miglioramento generale delle prestazioni del team rispetto allo sprint precedente, sebbene persistano alcune criticità nella stima del carico di lavoro e nella comunicazione interna. Alla luce di queste considerazioni, proponiamo le seguenti azioni per il prossimo sprint:
1. effettuare piu' pair-programming per risolvere problemi di comunicazione. Questo permetterebbe di mettere in diretta relazione persone di diversi ruoli e di far capire meglio il lavoro che si sta svolgendo.
2. Suddividere il carico di lavoro includendo anche PO e SM in attività di sviluppo, per massimizzare e parallelizzare il lavoro. Questo non sottraendo tempo alle attività già previste dai loro ruoli.


## Sprint 3
Valutazioni: 😃 😐 😡

| Cards            | Diego PO | Lorenzo SM | Samuele D | Fabio D | Emanuele D | Omar D |
| ---------------- | -------- | ---------- | --------- | --------| ---------- | ------ |
| Scrum Master     | 😃       | 😃         |  😃      | 😃       |  😃        | 😃     |
| Product Owner    | 😃       | 😃         |  😃      | 😃       |  😃        | 😃     |
| Development Team | 😃       | 😃         |  😃      | 😃       |  😃        | 😃     |
| Product Backlog  | 😃       | 😃         |  😃      | 😃       |  😃        | 😃     |
| Sprint Planning  | 😃       | 😃         |  😃      | 😃       |  😃        | 😃     |
| Sprint Goal      | 😃       | 😃         |  😃      | 😃       |  😃        | 😃     |
| Self Menagement  | 😃       | 😃         |  😃      | 😃       |  😐        | 😃     |

### Commenti personali
#### Product Owner - Diego Barbieri
Ottimo sprint, in assoluto il migliore. Nonostante un primo momento di insicurezza generale sul superamento delle task, ci siamo trovato a meta' settimna con tutte le us quasi completate: e' risultato necessario un "daily scrum on steroids" per decidere ulteriori us da terminare. Alla fine sono molto soddisfatto del risultato ottenuto, il tutto grazie a un codice ben organizzato in componenti.

#### Scrum Master - Lorenzo Peronese
Sprint eccellente, con il completamento di tutte le user stories assegnate. Gli imprevisti sono stati pochi e ben gestiti, grazie a uno sviluppo costante supportato dai tre daily scrum settimanali e da una comunicazione efficace, sia in chat che di persona in laboratorio. Il pair programming si è rivelato un grande successo, consentendo di risolvere molte problematiche in modo rapido ed efficace.

#### Developer - Omar Ayache
Sprint eccezionale, con obiettivi quasi completamente raggiunti. Diversi meeting scrum dinamici hanno permesso di ridefinire ulteriori attività. Il team ha dimostrato una significativa crescita, padroneggiando con competenza gli strumenti e comprendendo la loro complessità di utilizzo. L'organizzazione del codice in componenti ben strutturate ha contribuito al successo complessivo.

#### Developer - Samuele Musiani
Sprint decisamente positivo sotto molteplici aspetti. Il team di sviluppo ha dimostrato un'ottima cooperazione che ha portato al completamento di quasi tutte le US con un anticipo notevole rispetto alle tempistiche preventivate. Il team ha inoltre dimostrato una completa padronanza degli strumenti utilizzati, potendo quindi focalizzare le proprie risorse principalmente sulle attività di implementazione, ottimizzando così l'efficienza dell'intero processo di sviluppo. Questo sprint testimonia la maturità del team e pone solide basi per il proseguimento del progetto.

#### Developer - Emanuele Argonni
Sicuramente il migliore sprint fino ad ora, siamo riusciti a completare quasi tutte le task assegnate, nonostante avessimo un importante debito derivato dagli sprint precedenti.
Il team dei developer ha lavorato bene e si è dimostrato molto collaborativo. Il lavoro di squadra è stato fondamentale per il raggiungimento degli obiettivi.
Il pair programming è stato molto utile per l'integrazione backend-frontend e per risolvere i problemi che si sono presentati durante lo sprint.
Terminate le task assegnate, ci siamo dedicati al refactor del codice per migliorare la chiarezza e la manutenibilità del codice.
