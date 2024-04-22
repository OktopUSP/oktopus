<p align="center">
<img src="https://github.com/OktopUSP/oktopus/assets/83298718/fc05c512-951d-448c-8c31-1e0881783460"/></p>
<br/>
<ul>
    <li>
        <h4>Introduction:</h4>
    </li>
</ul>        
<p>
This repository aims to promote the development of a multi-vendor management platform for CPEs and IoTs. Any device that follows the TR-369 protocol can be managed. The main objective is to facilitate and unify device management, which generates countless benefits for the end user and service providers, suppressing the demands that today's technologies require: device interconnection, data collection, speed, availability and more.
</p>

<ul><li><h4>Sponsors:</h4></li></ul>

<a href="https://www.made4it.com.br/" target="_blank"><img src="https://github.com/OktopUSP/oktopus/assets/83298718/10da316f-e3fd-4e8d-93d3-5868cc2724e3" width="200px"/></a>

<ul><li><h4>Companies that use Oktopus:</h4></li></ul>

<a href="https://www.inango.com/" target="_blank"><img src="https://github.com/OktopUSP/oktopus/assets/83298718/3b3e65d9-33fa-46c4-8b24-f9e2a84a04a6" width="100px"/></a>

<p>If you'd like to know how to donate above <a href="https://github.com/sponsors/leandrofars">Github Sponsors</a> values, start a partnership or somehow to contribute to the project, email <a href="">leandro@oktopus.app.br</a>, every contribution is welcome, and the resources will help the project to move on. Also, if your company uses this project and you'd like your logo to appear up here, contact us.

--------------------------------------------------------------------------------------------------------------------------------------------------------

<ul>
    <li>
        <h4>💼 Commercial Support:</h4>
        <p>
            Our solution has an open-source software license, meaning you can modify/study the code and use it for free. You can perform all the configurations, allocate servers, and set it up on your network with the classic "do it yourself" approach, or save time and money: contact us for a quote and get commercial support.
        </p>
        <ul>
            <li> 
            Software customization according to your needs and preferences
            </li>
            <li> 
            Full support to your team
            </li>
            <li> 
            Affordable prices for companies of all sizes
            </li>
            <li> 
            Trust and assistance from experts
            </li>
            <li> 
            Complete solution for production environments, from the server provisionning to devices connection
            </li>
        </ul>
        <p>Contact <a href="">leandro@oktopus.app.br</a> via email and get a quote.</p>
    </li> 
</ul>

--------------------------------------------------------------------------------------------------------------------------------------------------------

<ul><li><h4>Infrastructure:</h4></li></ul>

![image](https://github.com/OktopUSP/oktopus/assets/83298718/aa6feb7f-ac32-465c-b166-aa6b3ee5b68a)

<ul>
    <li>
        <h4>API:</h4>
        <ul>
            <li> 
            <a href="https://documenter.getpostman.com/view/18932104/2s93eR3vQY#10c46751-ede9-4ea1-8ea4-264ebf539e5e">Documentation </a>
            </li>
        </ul>
    </li>
</ul>
<ul>
    <li>
        <h4>Quick start:</h4>
Run app using <u><b>Docker Compose</b></u>:
<pre>
user@user-laptop:~$ cd oktopus/deploy/compose
user@user-laptop:~/oktopus/deploy/compose$ COMPOSE_PROFILES=nats,controller,mqtt,stomp,ws,adapter,frontend,portainer docker compose up -d
</pre>
Oktopus deployment in <u><b>Kubernetes</b></u> still is in beta phase: <a href="https://github.com/OktopUSP/oktopus/blob/main/deploy/kubernetes/README.md"> Instructions for Kubernetes deployment</a><p></p>
        UI will open at port 3000:
        <img src="https://github.com/OktopUSP/oktopus/assets/83298718/65f7c5b9-c08d-479a-8a13-fdfc634b5ca2"/>

</li>
    <li>
      <h4>Device test agent (obuspa):</h4>
        <p>Run MQTT agent:</p>
        <pre>user@user-laptop:~/oktopus$ docker run -d -v $(pwd)/agent/oktopus-mqtt-obuspa.txt:/obuspa/oktopus-mqtt-obuspa.txt --network host --name obuspa-mqtt oktopusp/obuspa:latest obuspa -r /obuspa/oktopus-mqtt-obuspa.txt -p -v4 -i lo</pre>
        <p>Run Websockets agent:</p>
        <pre>user@user-laptop:~/oktopus$ docker run -d -v $(pwd)/agent/oktopus-websockets-obuspa.txt:/obuspa/oktopus-websockets-obuspa.txt --network host --name obuspa-websockets oktopusp/obuspa:latest obuspa -r /obuspa/oktopus-websockets-obuspa.txt -p -v4 -i lo</pre>
        <img src="https://github.com/OktopUSP/oktopus/assets/83298718/4599d566-eada-4313-8ae1-31dae82391de"/>
        <img src="https://github.com/OktopUSP/oktopus/assets/83298718/501b4ccd-6147-4957-9096-695134e34b5e"/>
    </li>
</ul>

--------------------------------------------------------------------------------------------------------------------------------------------------------

<ul>
    <li>
        <h4>Roadmap:</h4>
        <p>
            The project goals are organized with milestones that have a due date, just like a sprint. Those issues grouped in milestones are done and have their status updated in a kanban board.
        </p>
        <ul>
            <li> 
            <a href="https://github.com/OktopUSP/oktopus/milestones">Milestones </a>
            </li>
            <li> 
            <a href="https://github.com/orgs/OktopUSP/projects/1/views/2">Kanban Board </a>
            </li>
        </ul>
    </li>
</ul>

--------------------------------------------------------------------------------------------------------------------------------------------------------

<p>Are you going to use our project in your company? would like to talk about TR-369 and IoT management, we're online on <a href="https://join.slack.com/t/oktopustr-369/shared_invite/zt-1znmrbr52-3AXgOlSeQTPQW8_Qhn3C4g">Slack</a>.</p>

--------------------------------------------------------------------------------------------------------------------------------------------------------
<ul>
    <li>
        <h4>TR-069 ---> TR-369 :</h4>
    </li>
</ul>  
<p>
The advent of the Internet of Things brings countless opportunities and challenges for service providers, with over a billion devices across the globe today making use of <a href="https://www.broadband-forum.org/download /TR-069_Amendment-2.pdf">TR-069</a>, what is the future of the protocol and what can we expect ahead?
</p>
<p>
The CWMP (CPE Wan Management Protocol), better known as TR-069, opened many doors for the ecosystem of providers, through which it is possible to deliver services with agility, which meet or exceed customer expectations, with proactive management and secure network, also bearing in mind the lower cost and greater efficiency for service providers.
</p>
<p>
With the rise of what we now call the smart home, the Internet of Things and the demand for increasingly interconnected and cloud-based environments, new demands and obstacles have emerged, opening the door to the creation of a new form of communication that meets the needs of current market needs.
</p>
<p>
There is a fierce race to monetize the IoT devices that are now part of the connected home and other environments. As a result, many companies are creating their own proprietary solutions; this is understandable given such pressure generated by the promise of Smart Home monetization. Unfortunately, these applications contribute to a poor ecosystem, where a provider ends up dependent and limited to a vertical solution, from a single vendor. This generates an <b>low competition environment (which leads to greater risks), less innovation, and the potential for very high cost solutions</b>.
</p>
<p>
The technologies behind Wi-Fi, device-to-device connectivity, the Smart Home and IoTs are constantly evolving and improving. It is important that when service providers look for a solution, they look for something that is "future proof", always thinking ahead.
</p>
<p>
Seeking to solve the challenges mentioned above, providers and manufacturers together developed the USP (User Services Platform), defined by the Broadband Forum's TR-369 standard, which is the natural evolution of the TR-069. <b>This new standard is designed to be flexible, secure, scalable and standardized to meet the demands of a connected world today, and in the future.</b>
</p>

<ul>
    <li>
        <h4>Companies/Institutions involved in the creation of the TR-369:</h4>
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
        <h4>Topology:</h4>
    </li>
</ul>  

<img src="https://usp.technology/specification/architecture/usp_architecture.png"/>

![image](https://github.com/leandrofars/oktopus/assets/83298718/b1d5a0c7-4567-464c-bc9b-1956ef5c5f3b)

![image](https://github.com/leandrofars/oktopus/assets/83298718/7b46dc1f-5eb2-4a1b-8e77-376b0836948a)

<ul>
    <li>
        <h4>Protocols:</h4> 
        
![image](https://github.com/leandrofars/oktopus/assets/83298718/9b789f0b-cb0c-4cec-8b8e-767ba21bffae)
    </li>
</ul>

<ul>
    <li>
        <h4>Notifications/Data collection:</h4> 
        You can create notifications that fire on a value change, object creation and removal, complete operation, or an event.
        
![image](https://github.com/leandrofars/oktopus/assets/83298718/184899a3-52e7-491a-8ee7-7b442fe50719)
    </li>
</ul>

<ul>
    <li>
        <h4><a href="https://github.com/BroadbandForum/obuspa">OB-USP-A</a> (Open Broadband User Services Platfrom Agent):</h4> 
        <ul>
             <li>
             Designed for embedded software (~400kb on ARM)
             </li>
             <li>
             Encoded in C
             </li>
             <li>
             License <a href="https://opensource.org/license/bsd-3-clause/">BSD 3-Clause</a>
             </li>
             <li>
             Made for Linux environments
             </li>
        </ul>
    </li>
</ul>

<ul>
    <li>
        <h4>Data Analysis</h4>
The protocol has a mechanism called "Bulk Data", where it is possible to collect large volumes of data from the device, the data can be collected by HTTP, or another telemetry MTP defined in the TR standard, this data can be in JSON, CSV format or XML. This generates the opportunity to use AI on top of this data, obtaining relevant information that can be used for different purposes, from predicting events, KPIs, information for the commercial area, but also for the best configuration of a device.
    </li>
</ul>

<ul>
    <li>
        <h4>WiFi:</h4>
It has over 130 Wi-Fi configuration and diagnostics metrics, many of these settings and parameters are a trade-off between signal coverage area, latency and throughput. When deploying Wi-Fi systems, there is a tendency to maintain the same configuration on all clients, causing the technology to perform below expectations. Machine Learning combined with the data analysis mentioned in the previous topic makes it possible to automate the management and optimization of Wireless networks, where a big data approach is able to find the ideal configuration for each device.
        
![image](https://github.com/leandrofars/oktopus/assets/83298718/3d6fe3e8-3ca2-460b-9583-da89b42753f8)
    </li>
</ul>

<ul>
    <li>
       <h4>Commands:</h4>
         It is possible to perform commands remotely on the product, such as: firmware update, reboot, reset, search for neighboring networks, backup, ping, network diagnostics and many others.
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
<h4>Software Modules:</h4>
Currently, telecommunications giants and startups, publishing new software daily, slow delivery cycles and manual and time-consuming quality assurance processes make it difficult for integrators and service providers to compete. USP "Software Module Management" allows a containerized approach to the development of software for embedded devices, making it possible to drastically reduce the chance of error in software updates, it also facilitates the integration of third parties in a device, still keeping the firmware part isolated from Vendor.
<br/>
<img src="https://github.com/leandrofars/oktopus/assets/83298718/64664b0e-81cd-4a29-bbc5-b4186a04dfa2" width="50%"/>
    </li>
</ul>

--------------------------------------------------------------------------------------------------------------------------------------------------------

<p>Bibliographic sources: <a href="https://www.broadband-forum.org/download/MU-461.pdf">MU-461.pdf</a>, <a href="https:/ /usp.technology/specification/index.htm">TR-369.html</a>, <a href="https://drive.google.com/drive/folders/1N7FqK0PkDhjCN5s3OhQ_wmz9UcTSwRCX">USP Training Session Slides</usp.technology/specification/index.htm">TR-369.html</a></p>
