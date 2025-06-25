This project follows the following coding guidelines:
* You will explain decisions and design choices before implementing them.
* The architecture combines the best practices from Domain-Driven Design (DDD), Clean Architecture, and Hexagonal Architecture, creating a hybrid model that leverages the strengths of each.  All code must adhere to these designs.
* The codebase is written in Go, following best practices and conventions.
* Any function requiring a context parameter must have a context.Context parameter as the first argument and use ContextLogger to log messages.
* You must use Testify for assertions in tests.
* You must use GoMock for mocking in tests.
* You must not implement testing for the /examples folder.
* Implement unit tests for each package and when a new Go file is added.  The target is a minimum of 80% coverage of statements.  
* Unit tests should be fast, must exist in the same folder as the file being tested, and isolated from each other.
* Implement integration tests for each package as needed to adhere to best practices, with a minimum of 80% coverage of statements.
* Always use exhaustive error detection and handling, including retries and timeouts, where appropriate.
* Always update .svg files when updating .puml files.
* Always update the relevant.puml files whenever the architecture changes or code changes.
* Always update the documentation whenever there are relevant architecture changes or code changes.  This includes /docs, and any relevant .md files throughout the packages in the codebase.
* Do not embed examples in the documentation.  Rather, link to them in /examples.
* All new or existing README.md files must adhere to the COMPONENT_README_TEMPLATE.md structure.
* Every package must have a README.md file, even if it is empty.
