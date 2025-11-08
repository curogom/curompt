# Security & Privacy

> 🇰🇷 **Korean**: [한국어 버전](./SECURITY.ko.md)

- Default operation is **local-only**. No external network transmission.
- Captured data is stored in `~/.curompt`. Use `--no-store` if needed.
- Secret masking: Regex-based identification → `[REDACTED]` replacement.
- For dynamic evaluation (optional), a flag is provided to prevent storing sample outputs: `--no-samples-store`.
- Default log level is `info`. Secret logs are prohibited.
- Team mode (future): Model calls only through internal proxy.
