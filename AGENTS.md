# AI Agent Rules

This file defines the global rules that all AI agents must follow when working on this project.

## Rules

1. **Investigate before changing**
   - Identify and explain the root cause before making code changes.
   - Do not implement fixes based on assumptions.

2. **Keep changes minimal**
   - Make the smallest change necessary to solve the task.
   - Do not modify unrelated code or perform unnecessary refactoring.

3. **Preserve existing architecture**
   - Follow the existing project architecture, coding style, and conventions.
   - Do not introduce new dependencies or patterns unless explicitly requested.

4. **Never edit generated files**
   - Never edit `*_templ.go` files directly.
   - Treat them as generated code.

5. **Modify `.templ` files only**
   - Make UI changes only in the corresponding `.templ` files.
   - After editing, run `templ generate` to regenerate `*_templ.go`.

6. **Use project documentation when needed**
   - Read files under `context/` only when they are relevant to the current task or explicitly referenced by the user.

7. **Verify before finishing**
   - Ensure generated files are up to date.
   - Ensure the project builds successfully before completing the task.
   
8. **Ask instead of guessing**
   - If requirements are ambiguous, ask for clarification instead of making assumptions.

9. **Use PowerShell**
   - Always use Windows PowerShell-compatible commands. Never use CMD syntax.