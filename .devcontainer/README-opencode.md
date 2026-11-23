# OpenCode in GitHub Codespaces

This branch includes a browser-based OpenCode environment. The configuration is isolated to this branch and does not change the `main` branch or the Railway deployment.

## Start

1. Create a Codespace from this branch.
2. Wait for the container setup to finish. OpenCode is installed automatically and its web server starts on port `4096`.
3. In the Codespaces **Ports** panel, open the private forwarded port named **OpenCode Web**.
4. When prompted by the browser, use:
   - Username: `opencode`
   - Password: run `bash .devcontainer/show-opencode-password.sh` in the Codespaces terminal.

The password is generated inside the Codespace and is not stored in this repository.

## Use

OpenCode will open in the repository workspace. Choose and authenticate with an AI provider from the OpenCode interface before starting a coding session. Provider/API usage may incur separate charges.

To restart the web server manually:

```bash
bash .devcontainer/start-opencode.sh
```

To inspect the server log:

```bash
tail -n 100 ~/.config/opencode-codespaces/server.log
```

Port 4096 is configured as **private**, so it remains protected by GitHub Codespaces access controls.
