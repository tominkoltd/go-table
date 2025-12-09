# Changelog

## v0.3 — Live Redraw + Pointer Updates (2025-12-09)

### Added
- **Draw(true)** redraw mode  
  - Uses ANSI `ESC[{n}A` to move the cursor up by the number of lines drawn previously  
  - Enables flicker-free in-place table refreshing  
  - Useful for dashboards, counters, and status displays

- **Pointer-based Cell updates**
  - Accepts `*Cell` in `Push(...)`  
  - Mutating `cell.Data` (or color) reflects automatically on the next draw  
  - Perfect for simple terminal “animations” and live metric updates

### Changed
- Internal tracking of rendered line count (`drawnLines`) to support redraw mode  
- Minor internal cleanup around buffer building and ANSI composition

### Notes
First version explicitly targeting long-running CLI tools and process monitors.

---

## v0.2 — Width & Unicode Improvements

- Improved Unicode-safe rune handling  
- Better field width calculation for mixed characters  
- Border colorization refinements

---

## v0.1 — Initial Release

- Basic table rendering  
- ANSI colors, prefixes and suffixes  
- Flex and width-based layout system
