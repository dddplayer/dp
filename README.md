# ddd-player

## Goal

Make DDD model a workable model

## How

* Solve problems mentioned by [Domain-Driven Design Reference](https://www.domainlanguage.com/wp-content/uploads/2016/05/DDD_Reference_2015-03.pdf)
* Provide all process modelling supporting of [DDD Starter Modelling Process](https://github.com/ddd-crew/ddd-starter-modelling-process)
* [Stack Overflow](https://stackoverflow.com/questions/tagged/domain-driven-design%20go?sort=Newest&edited=true)

## Reference

* [Martin Fowler](https://martinfowler.com/)
* [Eric Evans](https://www.domainlanguage.com/)
* [DDD Crew](https://github.com/ddd-crew)


Why:

* Ubiquitous language
* Focus on Core Domain
* Collaboration
* Continuous

## Guidance

- **Use Cases Driven and stick to ddd-crew and domain language**
- Tell story by using our **domain message flow**
- **Domain flow** to help conversation between domain expert and delivery team
- Domain diff to help conversation between tech manager (TL/TP) and delivery team
- Event storming ?
- Context Map ?
- Modulation or registry model?

## Expectations

* Code: source code analysis
* Arch: Graph creator and output
  * One repo for objects, one repo for relations
  * Build base graph with objects and relations
  * Output different type of Graphs
* Domain: Support for Arch with domain information

Todo:

* Automatically recognize domain from code - Done
* Domain model concepts not appear in aggregate Domain Model, need to add it on - Done 
  * Domain message flow - ubiquitous language for communication / story telling
    * How many services in this application, and how domain support those services?
    * What's the core domain - except application/infrastructure/interface, others are domains
    * What's the information we have in the flow - domain event?
    * How message flow looks like?
    * How many domains are there?
    * Domain
      * How many aggregates inside this domain, and how they linked with each other - Done
      * Implementation details of domain - Doing
    * Services
      * Dot Service
        * Strategic module Dot Graph
        * Tactic module Dot Graph
        * Dot Graph
          * Build
          * String
          * Append node
      * Code
        * Language
          * Golang
            * Parser - AST/SSA/Call Graph - internal concept
        * Code Structure
      * Domain Model Service
        * Core Domain ?
        * Get Domain
        * Add domain component - internal concept
        * Architecture (package name, path, filename, objects)
          * Objects
          * Code Visitor
            * Node handler
            * Link handler
          * Hexagon - Done
            * Application
              * Service
            * Interface
            * Infrastructure
              * Message system
          * Plain
        * Domain
          * Input from Repository
          * Build
          * Name
          * Aggregates
          
* Dot node and edge with attributes management, template focus on structure, make template org automatically
* Save SVG, share SVG
* Convert Markdown to DDD plus hexagon architecture - Done
* Refactor Markdown based on DDD Player analysing result

```text
project/
├── cmd/
│   └── main.go
├── internal/
│   ├── application/
│   │   └── ...
│   ├── domain/
│   │   ├── aggregate/
│   │   ├── entity/
│   │   ├── repository/
│   │   ├── service/
│   │   ├── valueobject/
│   │   └── ...
│   ├── infrastructure/
│   │   ├── persistence/
│   │   ├── messaging/
│   │   ├── external/
│   │   └── ...
│   └── interfaces/
│       ├── rest/
│       ├── grpc/
│       ├── eventhandler/
│       ├── job/
│       └── ...
└── pkg/
    ├── utils/
    ├── config/
    └── ...

```

* Aggregate root - next (Domain -> Model -> Factory -> Aggregate root) - Definition

- DDD 战略图：战略图描述领域的上下文边界、限界上下文以及它们之间的关系。使用合适的布局来表达领域的整体结构。
- DDD 战术图：战术图描述领域的设计模式、聚合、实体、仓储等元素，以及它们之间的关系。使用适当的关系标记来表示它们之间的关系。

* User journey animation - story telling
* Domain Model Management
  * Changes link with ADR and commit log
  * Changes notification - ADR generated - different role touch the same code
* Separate message with type, which will help us to focus on specific message

Goal: Responsive to the user's need.

* Support Repository domain component
* LoadAllSyntax deprecated
* tests coverage improvement > 80%
* Execution animation
* Code link


* Call graph analysis algo options - Analysis performance improvement
  * RA, reachability analysis: graph = static.CallGraph(prog)
  * CHA, class hierarchy analysis: graph = cha.CallGraph(prog)
  * RTA, rapid type analysis: graph = rta.Analyze(roots, true).CallGraph
  * Pointer - current one

## Bounded Context

Business Capabilities

我最推荐的方法是结合持续集成和自动化测试以及代码审查。

持续集成和自动化测试可以在代码提交前自动运行测试套件，确保代码的正确性和稳定性。通过编写全面的测试用例并监控测试覆盖率，可以及早发现代码与战略图的不匹配问题。

另外，代码审查是一种重要的实践，可以通过团队的协作来发现和解决代码与战略图之间的不一致情况。代码审查能够提供更深入的理解和讨论，确保团队成员对代码和战略图的一致性有共识。

通过持续集成和自动化测试，以及代码审查的组合使用，可以在开发过程中持续地监测代码和战略图的匹配度。这样可以及时发现潜在的不一致问题，并促使团队成员相互之间进行反馈和讨论，从而保持代码和战略图的一致性，提高软件的质量和可维护性。


当发现代码与战略图不符时，可以自动生成ADR（Architecture Decision Record）并发送通知给相关团队成员。ADR是一种记录架构决策和背后原因的文档，它可以用于记录战略图的变更和相关的讨论。

以下是一个可能的流程：

检测不符：通过代码审查、静态代码分析工具、持续集成或自动化测试，发现代码与战略图不符的情况。

自动生成ADR：创建一个自动化脚本或工具，可以根据检测到的不符情况自动生成ADR文档。ADR应包含以下内容：

描述不符的具体问题和与战略图的差异。
解释不符的原因或背后的决策。
提供解决方案或建议以使代码与战略图保持一致。

发送通知：一旦ADR生成，可以发送通知给相关的团队成员，如开发人员、架构师、产品经理等。通知可以通过电子邮件、即时消息工具或项目管理工具进行发送，以确保相关人员能够及时了解到不符情况。

讨论和决策：在ADR通知中，鼓励团队成员参与讨论和决策，以确定如何处理不符情况。讨论可以在代码审查会议、团队会议或专门的架构讨论会上进行。

更新和修正：根据讨论和决策的结果，相应的更新和修正代码或战略图。确保代码和战略图保持一致，并在ADR中记录任何更改或更新。

通过自动生成ADR并发送通知，可以提高团队对代码和战略图一致性的关注，并促使团队成员积极参与解决问题。这样可以确保及时记录和处理不符情况，以保持代码的质量和架构的一致性。