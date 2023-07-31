<p align="center">
<img src="https://user-images.githubusercontent.com/83298718/220207485-8c2aac78-95eb-4b43-b23e-c4bfa6cd30e6.png"/>
</p>
<br/>
<ul>
    <li>
        <a href="https://github.com/leandrofars/oktopus/blob/main/README.en.md">Readme in English</a>
    </li>
</ul>
<ul>
    <li>
        <h4>Introdução:</h4>
    </li>
</ul>        
<p>
Este repositório tem como intuito fomentar o desenvolvimento de uma plataforma de gerência multi-vendor para IoTs. Todo dispositivo que seguir o protocolo TR-369 poderá ser gerenciado. O objetivo principal é facilitar e unificar a gerência de dispositivos, o que gera inúmeros benefícios para o usuário final e prestadores de serviços, suprimindo as demandas que as tecnologias de hoje exigem: interconexão de dispositivos, coleta de dados, rápidez, disponibilidade e muito mais.
</p>
<ul>
    <li>
        <h4>TR-069 ---> TR-369 :</h4>
    </li>
</ul>  
<p>
O advento da Internet das Coisas traz inúmeras oportunidades e desafios pra os prestadores de serviços, com mais de um bilhão de dispositivos espalhados pelo globo hoje, fazendo uso do <a href="https://www.broadband-forum.org/download/TR-069_Amendment-2.pdf">TR-069</a>, qual é o futuro do protocolo e o que podemos esperar pela frente?
</p>
<p>
O CWMP(CPE Wan Management Protocol), mais conhecido como TR-069, abriu muitas portas para o ecossistema de provedores, por meio dele é possível entregar serviços com agilidade, que servem ou ultrapassam as expectativas do cliente, fazendo uma gestão pró-ativa e segura da rede, tendo em vista também o menor custo e a maior eficiência para os prestadores de serviços.
</p>
<p>
Com a ascensão do que hoje chamamos de casa inteligente, a Internet das Coisas e a demanda por ambientes cada vez mais interconectadas e baseados em nuvem, novas demandas e obstáculos surgiram, abrindo a porta para a criação de uma nova forma de comunicação que supra as necessidades do mercado atual.
</p>
<p>
Existe uma corrida acirrada para monetizar os dispostivos IoT que hoje fazem parte da casa conectada e de outros ambientes. Como resultado disso, muitas empresas estão criando suas próprias soluções proprietárias; isso é compreensível dada tamanha pressão gerada pela promessa da monetização da Casa Inteligente. Infelizmente, essas aplicações contribuem para um ecossistema pobre, onde um provedor acaba dependente e limitado a uma solução vertical, de um único Vendor. Isso gera um <b>ambiente de pouca competição (o que leva a maiores riscos), menos inovação, e o potencial de soluções com custos muito elevados</b>. 
</p>
<p>
As tecnologias por trás do Wi-Fi, a conectividade entre dispositivos, a Casa Inteligente e os IoTs estão em constante evolução e aprimoramento. É importante que quando os prestadores de serviços forem buscar uma solução, busquem por algo que seja a "prova de futuro", pensando sempre adiante.
</p>
<p>
Buscando resolver os desafios citados anteriormente, provedores e fabricantes juntos, desenvolveram o USP (User Services Platform), definido pela norma TR-369 da Broadband Forum, sendo este, a evolução natural do TR-069. <b>Este novo padrão foi desenhado para ser flexível, seguro, escalonável e padronizado, para atender as demandas de um mundo conectado hoje, e no futuro.</b>
</p>

<ul>
    <li>
        <h4>Empresas/Instituições envolvidas na criação do TR-369:</h4>
        <ul>
            <li> 
            Google
            </li>
            <li> 
            Nokia
            </li>
            <li> 
            Huawei
            </li>
            <li> 
            Axiros
            </li>
            <li> 
            Orange
            </li>
            <li> 
            Commscope
            </li>
            <li> 
            Assia
            </li>
            <li> 
            AT&AT
            </li>
            <li> 
            NEC
            </li>
            <li> 
            Arris
            </li>
            <li> 
            QA Cafe
            </li>
        </ul>
    </li>
</ul> 

--------------------------------------------------------------------------------------------------------------------------------------------------------

<ul>
    <li>
        <h4>Topologia:</h4>
    </li>
</ul>  

<img src="https://usp.technology/specification/architecture/usp_architecture.png"/>

![image](https://github.com/leandrofars/oktopus/assets/83298718/b1d5a0c7-4567-464c-bc9b-1956ef5c5f3b)

![image](https://github.com/leandrofars/oktopus/assets/83298718/7b46dc1f-5eb2-4a1b-8e77-376b0836948a)

<ul>
    <li>
        <h4>Protocolos:</h4> 
        
![image](https://github.com/leandrofars/oktopus/assets/83298718/9b789f0b-cb0c-4cec-8b8e-767ba21bffae)
    </li>
</ul>

<ul>
    <li>
        <h4>Notificações/Coleta de dados:</h4> 
        É possível criar notificações que são disparadas em uma mudaça de valor, criação e remoção de objeto, operação completa, ou um evento.
        
![image](https://github.com/leandrofars/oktopus/assets/83298718/184899a3-52e7-491a-8ee7-7b442fe50719)
    </li>
</ul>

<ul>
    <li>
        <h4><a href="https://github.com/BroadbandForum/obuspa">OB-USP-A</a> (Open Broadband User Services Platfrom Agent):</h4> 
        <ul>
            <li> 
            Desenhado para software embarcado (~400kb em ARM)
            </li>
            <li> 
            Codificado em C
            </li>
            <li> 
            Licença <a href="https://opensource.org/license/bsd-3-clause/">3-Clause BSD</a>
            </li>
            <li> 
            Feito para ambientes linux
            </li>
        </ul>
    </li>
</ul>

<ul>
    <li>
        <h4>Análise de Dados</h4> 
O protocolo possui um mecanismo chamado "Bulk Data", onde é possível recolher grandes volumes de dados do dispositivo, os dados podem ser recolhidos por HTTP, ou outro MTP de telemetria definido na norma do TR, esses dados podem estar em formato JSON, CSV ou XML. Isso gera a oportunidade de utilizar IA em cima desses dados, obtendo informações relevantes que podem ser usadas tendo diversas intenções, desde a predição de eventos, KPIs, informações para a área comercial, mas também para a melhor configuração de um dispositivo.
    </li>
</ul>

<ul>
    <li>
        <h4>Wi-Fi:</h4> 
Possui mais de 130 métricas de configuração e diagnóstico de Wi-Fi, muitas dessa configurações e parâmetros são uma troca entre área de cobertura do sinal, latência e throughput. Ao implantar sistemas Wi-Fi, tende-se a manter a mesma configuração em todos os clientes, fazendo com que a tecnologia tenha uma performance abaixo do esperado. A Machine Learning aliada à análise de dados citada no tópico anterior, torna possível automatizar o gerenciamento e a optimização de redes Wireless, onde uma abordagem de big data é capax de encontrar a configuração ideal para cada dispositivo.
        
![image](https://github.com/leandrofars/oktopus/assets/83298718/3d6fe3e8-3ca2-460b-9583-da89b42753f8)
    </li>
</ul>

<ul>
    <li>
        <h4>Comandos:</h4>
        É possível realizar comandos remotamente no produto, como por exemplo: atualização de firmware, reboot, reset, busca de redes vizinhas, backup, ping, diagnósticos de rede e muitos outros.
    </li>
</ul>

<ul>
    <li>
        <h4>IoT:</h4>
<div align="center">
<img src="https://github.com/leandrofars/oktopus/assets/83298718/a2a12d9d-05a0-428b-ba3f-1ad83c876301" width="90%"/>
<br/>
<img src="https://github.com/leandrofars/oktopus/assets/83298718/91a87f43-3de7-42bd-a689-a4e14eecf5c0" width="60%"/>
<br/>
<img src="https://github.com/leandrofars/oktopus/assets/83298718/73e2e360-d53e-494e-9a50-60c83dae75df" width="60%"/>
<div>
    </li>
</ul>

<ul>
    <li>
        <h4>Módulos de Software:</h4>
Atualmente, gigantes das telecomunicações e startups, publicam software novo diariamente, ciclos de entrega lentos e processos de garantia de qualidade manuais e demorados, torna difícil a competição para integradores e prestadores de serviços. USP "Software Module Management" permite uma abordagem containerizada ao desenvolvimento de software de dispositivos embarcados, tornando possível diminuir drasticamente a chance de erro em atualições de software, também facilita a integração de terceiros em um dispositvo, manténdo ainda assim, isolada a parte de firmware do Vendor. 
<br/>
<img src="https://github.com/leandrofars/oktopus/assets/83298718/64664b0e-81cd-4a29-bbc5-b4186a04dfa2" width="50%"/>
    </li>
</ul>



--------------------------------------------------------------------------------------------------------------------------------------------------------

<ul><li><h4>Infraestrutura:</h4></li></ul>

![image](https://github.com/leandrofars/oktopus/assets/83298718/aa22edbb-bc82-4330-9678-650011bce5a8)

<ul>
    <li>
        <h4>API:</h4>
        <ul>
            <li> 
            <a href="https://documenter.getpostman.com/view/18932104/2s93eR3vQY#10c46751-ede9-4ea1-8ea4-264ebf539e5e">Documentação </a>
            </li>
            <li> 
            <a href="https://www.postman.com/docking-module-astronomer-46169629/workspace/oktopus">Workspace de testes e desenvolvimento</a>
            </li>
        </ul>
    </li>
</ul>
<ul>
     <li>
         <h4>Desenvolvedor:</h4>
Execute o aplicativo usando o Docker:
<pre>
leandro@leandro-laptop:~$ cd oktopus/devops
leandro@leandro-laptop:~/oktopus/devops$ docker run
</pre>
     </li>
     <li>
Compilação manual básica e execução:
     <ul>
         <li>
         <b>Dependências:</b> Versão node: v14.20.0 | Versão Go: v1.18.1
         </li>
         <li>
         Broker MQTT:
             <pre>leandro@leandro-laptop:~$ vá executar oktopus/backend/services/mochi/cmd/main.go -redis "127.0.0.1:6379"</pre>
         </li>
         <li>
         Controlador TR-369:
             <pre>
leandro@leandro-laptop:~$ execute oktopus/backend/services/controller/cmd/oktopus/main.go -u root -P root -mongo "mongodb://127.0.0.1:27017"</pre>
         </li>
         <li>
         Servidor socketio:
             <pre>
leandro@leandro-laptop:~$ cd oktopus/backend/services/socketio && npm i && npm start</pre>
         </li>
         <li>
         Front-end:
             <pre>
leandro@leandro-laptop:~$ cd oktopus/frontend && npm i && npm run dev</pre>
         </li>
     </ul>
</li>
OBS: Não use essas instruções em produção. Para implementar o projeto em produção você pode usar mais recursos que já estão disponíveis no Oktopus, mas levaria mais tempo para explicar neste README. Em breve, haverá mais ajuda e explicações sobre essas configurações extras necessárias.
<li>
Configurações dispositivo de teste:
    </li>
<ul>
        <li>
<b>Device.LocalAgent.</b>
<pre>
"CertificateNumberOfEntries": "0",
"Controller.1.Alias": "",
"Controller.1.AssignedRole": "",
"Controller.1.BootParameterNumberOfEntries": "0",
"Controller.1.ControllerCode": "",
"Controller.1.Enable": "true",
"Controller.1.EndpointID": "oktopusController",
"Controller.1.InheritedRole": "Device.LocalAgent.ControllerTrust.Role.1",
"Controller.1.MTP.1.Alias": "",
"Controller.1.MTP.1.Enable": "true",
"Controller.1.MTP.1.MQTT.Reference": "Device.MQTT.Client.1",
"Controller.1.MTP.1.MQTT.Topic": "oktopus/v1/controller",
"Controller.1.MTP.1.Protocol": "MQTT",
"Controller.1.MTPNumberOfEntries": "1",
"Controller.1.PeriodicNotifInterval": "15",
"Controller.1.PeriodicNotifTime": "0001-01-01T00:00:00Z",
"Controller.1.ProvisioningCode": "",
"Controller.1.USPNotifRetryIntervalMultiplier": "2000",
"Controller.1.USPNotifRetryMinimumWaitInterval": "5",
"ControllerNumberOfEntries": "1",
"ControllerTrust.ChallengeNumberOfEntries": "0",
"ControllerTrust.CredentialNumberOfEntries": "0",
"ControllerTrust.Role.1.Alias": "cpe-1",
"ControllerTrust.Role.1.Enable": "true",
"ControllerTrust.Role.1.Name": "Full Access",
"ControllerTrust.Role.1.Permission.1.Alias": "cpe-1",
"ControllerTrust.Role.1.Permission.1.CommandEvent": "r-xn",
"ControllerTrust.Role.1.Permission.1.Enable": "true",
"ControllerTrust.Role.1.Permission.1.InstantiatedObj": "rw-n",
"ControllerTrust.Role.1.Permission.1.Obj": "rw-n",
"ControllerTrust.Role.1.Permission.1.Order": "0",
"ControllerTrust.Role.1.Permission.1.Param": "rw-n",
"ControllerTrust.Role.1.Permission.1.Targets": "Device.",
"ControllerTrust.Role.1.PermissionNumberOfEntries": "1",
"ControllerTrust.RoleNumberOfEntries": "2",
"EndpointID": "os::test-000000000001",
"MTP.1.Alias": "",
"MTP.1.Enable": "false",
"MTP.1.MQTT.PublishQoS": "0",
"MTP.1.MQTT.Reference": "Device.MQTT.Client.1",
"MTP.1.MQTT.ResponseTopicConfigured": "oktopus/v1/controller",
"MTP.1.MQTT.ResponseTopicDiscovered": "oktopus/v1/agent/os::4851CF-000000000002",
"MTP.1.Protocol": "MQTT",
"MTP.1.Status": "Down",
"MTPNumberOfEntries": "1",
"RequestNumberOfEntries": "0",
"SoftwareVersion": "5.0.0",
"SubscriptionNumberOfEntries": "0",
"SupportedFingerprintAlgorithms": "SHA-1, SHA-224, SHA-256, SHA-384, SHA-512",
"SupportedProtocols": "STOMP, CoAP, MQTT, WebSocket",
"UpTime": "42"
</pre>
    </li>
    <li>
<b>Device.MQTT.Client.1</b>
    </li>
    <pre>
"Alias": "cpe-1",
"BrokerAddress": "10.100.250.4",
"BrokerPort": "1883",
"CleanSession": "false",
"CleanStart": "false",
"ClientID": "",
"ConnectRetryIntervalMultiplier": "2000",
"ConnectRetryMaxInterval": "30720",
"ConnectRetryTime": "5",
"Enable": "true",
"KeepAliveTime": "30",
"Name": "",
"Password": "",
"ProtocolVersion": "5.0",
"RequestProblemInfo": "false",
"RequestResponseInfo": "false",
"ResponseInformation": "oktopus/v1/agent/os::4851CF-000000000002",
"Status": "Connected",
"Subscription.1.Alias": "cpe-1",
"Subscription.1.Enable": "false",
"Subscription.1.QoS": "1",
"Subscription.1.Topic": "oktopus/v1/agent",
"SubscriptionNumberOfEntries": "1",
"TransportProtocol": "TCP/IP",
"Username": "test"</pre>
</ul>
</ul>

--------------------------------------------------------------------------------------------------------------------------------------------------------
<p>Vai usar nosso projeto na sua empresa? gostaria de conversar sobre o TR-369 e gerenciamento de IoTs, estamos online no <a href="https://join.slack.com/t/oktopustr-369/shared_invite/zt-1znmrbr52-3AXgOlSeQTPQW8_Qhn3C4g">Slack</a>.</p>
<p>Caso você tenha interesse em informações internas sobre o time e nossas pretensões acesse nossa <a href="https://github.com/leandrofars/oktopus/wiki">Wiki</a>.</p>

--------------------------------------------------------------------------------------------------------------------------------------------------------

<p>Fontes bibliográficas: <a href="https://www.broadband-forum.org/download/MU-461.pdf">MU-461.pdf</a>, <a href="https://usp.technology/specification/index.htm">TR-369.html</a>, <a href="https://drive.google.com/drive/folders/1N7FqK0PkDhjCN5s3OhQ_wmz9UcTSwRCX">USP Training Session Slides</a></p>

