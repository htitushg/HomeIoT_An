
![HomeNCo-logo](https://github.com/user-attachments/assets/2fd0196a-1147-4c46-a31b-b0088f8e85bf)

---

## Presentation

This project is made by Computer Science students for an assignment.

**Home&Co** _(Co for Connect, Control, Company)_ is a project using **ESP32** dev boards and a **Raspberry Pi 3B+** to monitor, control and automate home appliances like the lights, heaters, front door, water valves for watering plants, etc.

Its objectives are to provide **security**, **control**, **comfort**, **peace of mind** and **energy saving** to your home.

### 1. Security

Our front door device will control the access policy of the front door of your home:
- **RFID** sensor needing a **badge** to open the door
- **bell button** ringing the bell and **notifying you on your phone**
- **locking mechanism** controlled **remotely** (no need to give keys to anybody)
- **presence detector** to **monitor** the activity in front of the door
- **house locking function** to **automatically turn off** all the lights and use all presence detectors as **intrusion detectors**

### 2. Control

You can use the website to **remotely control at any time** the _lights_, the _door_, the _heaters_ and the _watering of plants_.

You have the **full control of your house** from any device.

### 3. Comfort

You can **automate the heaters** to activate at any time of the day, to find a warm house when coming back from a hike, vacation or any activity.

### 4. Peace of mind

You won't have to go through your check list three times or ask your neighbour to water your plants before going out!

You can **monitor your house from any device** using the website in **real time**: from the _presence detectors_ to the _lights_, _heaters_, _humidity of the plants_ and _temperature of any room_.

### 5. Energy saving

No need to keep heating your home when you're not there!

You'll be able to control and monitor the **heaters** remotely and according to your needs and presence at home.

You'll also be able to **monitor the energy of specific high consumption appliances** with our **smart power outlet** and control them remotely, and even give them schedules to be working or not.


## How it works

### Needs analysis

- Fully functional home with remote control (from phone or any other device)
- Security: no direct access to sensors or actuators (ESP32), to database, MQTT or backend application.
- Raspberry Pi 3 or more for WiFi access point, MQTT, backend app and webserver.
- ESP32 for sensors & actuators

### Technical choices

- Golang for the backend app and webserver with GORM: easy to use system language with good web libraries and ORM
- PostgreSQL for the database: the best Open Source database
- Mosquitto for the MQTT broker: good Open Source MQTT broker (up to date and good documentation)
- RaspAP for the WiFi Access Point in the Raspberry Pi: the most updated and documented Access Point software we found
- Websocket between JavaScript clients and Golang server: best solution for event driven actualization in the frontend.
- PlatformIO and C/C++ for the ESP32: the most manageable software to compile and upload C/C++ code on ESP32


### Architecture

```mermaid
stateDiagram-v2
    S: Raspberry Pi Server
    state S {
        MQTT: MQTT Broker
        A: Server Application
        DB: PostgreSQL Database
        
        A --> MQTT: Monitors & send commands
        MQTT --> A: Notify
        
        A --> DB: Updates the database
    }
    
    LD: Light
    state LD {
        App: Application
        B: Broker
        LC: Light Controller
        LS: Luminosity Sensor
        TS: Temperature Sensor
        PD: Presence Detector
        
        LC --> B: Subscribes /lightController channel
    }
    
    B --> MQTT: Connects
    
```

### Database

```mermaid
---
title: Entity Relationship Diagram
---
erDiagram
    DATA {
        int id
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
        string device_id
        int user_id
        int module_id
        string module_name
        string module_value
    }
    DEVICES {
        string id
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
        int location_id
        string type
        string name
    }
    LOCATIONS {
        int id
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
        string type
        string name
    }
    MODULES {
        int id
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
        string device_id
        string name
        string value
    }
    USERS {
        int id
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
        string name
        string password_hash
        string email
        string phone_number
        string status
    }
    SESSIONS {
        string token
        string data
        timestamp expiry
    }
    
    LOCATIONS ||--o{ DEVICES : contains
    DEVICES }|--|| MODULES : "is composed of"
    DATA ||--o{ DEVICES : "is about"
    DATA |o--o{ USERS : "owns"
```

### Incoming features

- Monitoring sensors with the web interface
- Controlling actuators from the web interface

### MQTT

```mermaid
---
title: MQTT topics/channels
---

flowchart LR
    system["`_SYSTEM NAME_`"]
    location_type["`_LOCATION TYPE_`"]
    location_id["`_LOCATION ID_`"]
    device_type["`_DEVICE TYPE_`"]
    device_id["`_DEVICE ID_`"]
    module["`_MODULE_`"]
    
    home["`**home**`"]
    
    kitchen["`**kitchen**`"]
    garden["`**garden**`"]
    room["`**room**`"]
    
    kitchen_id["`**1**`"]
    room1_id["`**2**`"]
    room2_id["`**3**`"]
    garden_id["`**4**`"]
    
    light["`**light**`"]
    water_plant["`**waterPlant**`"]
    
    device1["`**ESP32-af6f-43b1-a20e**`"]
    device2["`**ESP32-db87-47b7-998d**`"]
    device3["`**ESP32-808e-4c62-9c4f**`"]
    device4["`**ESP32-0158-4d61-85c2**`"]
    
    light_controller["`**lightController**`"]
    luminosity_sensor["`**luminositySensor**`"]
    water_valve["`**waterValve**`"]
    humidity_sensor["`**humiditySensor**`"]
    
    subgraph OUR SYSTEM NAME
        system
        home
    end
    subgraph LOCATION
        subgraph  
        location_type
        kitchen
        garden
        room
        end
        subgraph  
        location_id
        kitchen_id
        garden_id
        room1_id
        room2_id
        end
    end
    subgraph DEVICE 
       subgraph  
        device_type
        light
        water_plant
        end
        subgraph  
        device_id
        device1
        device2
        device3
        device4
        end 
    end
    subgraph SPECIFIC CHANNEL
        module
        light_controller
        luminosity_sensor
        water_valve
        humidity_sensor
    end

    system  --- location_type   --- location_id     --- device_type     --- device_id   --- module
    home    --- kitchen     --- kitchen_id  --- light           --- device1     --- light_controller & luminosity_sensor
    home    --- garden      --- garden_id   --- water_plant     --- device2     --- water_valve & humidity_sensor
    home    --- room        --- room1_id    --- light           --- device3     --- light_controller & luminosity_sensor
    home    --- room        --- room2_id    --- light           --- device4     --- light_controller & luminosity_sensor
```


### Class diagram

```mermaid
---
title: Device class diagram
---
classDiagram
    class Application {
        # Application *app$
        # String location
        # unsigned int locationID
        # String root_topic
        # bool wait_for_setup
        # unsigned int publish_interval
        # unsigned long lastPublishTime
        # WiFiClient network
        # Broker *broker
        # IModule *lightController
        # IModule *lightSensor
        # IModule *luminositySensor
        # IModule *presenceDetector
        # IModule *temperatureSensor
        # IModule *consumptionSensor
        
        # isWaitingForSetup() bool
        # onSetupMessage(char payload[]) void
        # reset() void
        # setupModule(const char* name, const char* value) void
        # messageHandler(MQTTClient *client, char topic[], char payload[], int length) void$
        # unsubscribeAllTopics() void
        # setRootTopic() void

        + getInstance() Application*$
        + Application() Application
        + brokerLoop() void
        + startup() void
        + init(WiFiClient wifi) void
        + sensorLoop() void$
    }
    class Broker {
        # MQTTClient mqtt
        # WiFiClient wifi
        # String root_topic

        # Broker(WiFiClient network) Broker

        + newBroker(WiFiClient network, void cb(MQTTClient *client, char topic[], char bytes[], int length)) Broker*$
        + sub(const String &module_name) void
        + pub(const String &module_name, const String &value) void
        + unsub(const String &module_name) void
        + setRootTopic(const String &topic) void
        + loop() void
    }
    class IModule {
        <<interface>>
        # Broker *broker
        # String name
        + setValue(const char * value) void*
        + getValue() const String*
        + getValueReference() const void**
        + getName() String
    }
    class ModuleFactory {
        + newModule(Broker *broker, String type) IModule*$
    }
    class IObservable {
        <<interface>>
        + ~IObservable()*
        + Attach(IObserver *observer) void*
        + Detach(IObserver *observer) void*
        + Notify() void*
    }
    class IObserver {
        <<interface>>
        + ~IObserver()*
        + Update(const String &value) void*
    }
    class LightController {
        # Broker *broker
        # String name
        # bool value
        
        + LightController(Broker *broker, bool value) LightController
        + setValue(const char * value) void
        + getValue() const String
        + getValueReference() const * void
        + getName() String
        + Attach(IObserver *observer) void
        + Detach(IObserver *observer) void
        + Notify() void
        + Update(const String &value) void
    }
    class PresenceDetector {
        # Broker *broker
        # String name
        # bool value

        + PresenceDetector(Broker *broker, bool value) PresenceDetector
        + setValue(const char * value) void
        + getValue() const String
        + getValueReference() const * void
        + getName() String
        + Attach(IObserver *observer) void
        + Detach(IObserver *observer) void
        + Notify() void
        + Update(const String &value) void
    }
    Application *.. IModule
    ModuleFactory ..> IModule
    Application --> ModuleFactory
    Application *.. Broker
    LightController ..|> IModule
    PresenceDetector ..|> IModule
    IObservable *.. IObserver
    IModule ..|> IObserver
    IModule ..|> IObservable
```

### Startup / Setup

#### Startup/Setup message
```json
{
  "id": "ESP32-af6f-43b1-a20e",
  "type": "light",
  "location_id": 3,
  "location_type": "room",
  "location_name": "room 3",
  "modules": [
    {
      "name": "lightController",
      "value": "False"
    },
    {
      "name": "lightSensor",
      "value": "True"
    },
    {
      "name": "luminositySensor",
      "value": "150.0"
    },
    {
      "name": "presenceDetector",
      "value": "True"
    },
    {
      "name": "temperatureSensor",
      "value": "22.5"
    },
    {
      "name": "consumptionSensor",
      "value": "32.45"
    }
  ]
}
```

#### Setup sequence diagram

```mermaid
sequenceDiagram
    autonumber
    box Teal Device
    participant S as Setup
    participant A as Application
    participant B as Broker
    participant LiC as LightController
    participant LiS as LightSensor
    participant PD as PresenceDetector
    participant LuS as LuminositySensor
    participant TS as TemperatureSensor
    participant CS as ConsumptionSensor
    end
    box Green Server
    participant MQTT as MQTT Broker
    participant Server as Go Application
    participant DB
    end
    S ->> A: Instanciate
    A ->> B: Instanciate
    B ->> MQTT: Connect
    A ->> B: Prepare `startup` message
    B ->> MQTT: Send `startup` message
    MQTT ->> Server: Relay `startup` message from Device
    Server ->> DB: Check if Device exists and creates it if necessary
    DB ->> Server: Send Device data back to prepare `setup` message
    Server ->> MQTT: Send `setup` message
    MQTT ->> B: Relay `setup` message
    B ->> A: Parse `setup` message
    A ->> A: Update `location` and `locationID`
    A ->> LiC: Update value
    A ->> LiS: Update value
    A ->> PD: Update value
    A ->> LuS: Update value
    A ->> TS: Update value
    A ->> CS: Update value
    A ->> S: Setup complete
```