# Chronicle

**Explore Céline Dion's life. Rebuild and play any subject's story with AI.**

From Québec to Eurovision and the world stage, follow Céline Dion's performances,
travels, collaborators, and milestones through a timeline you can play aloud.

Chronicle is a local academic research workspace that automatically investigates
a subject's documented history and rebuilds its story with sources you can inspect.
Built for students, educators, researchers, and curious minds.

![Céline Dion's researched history in Chronicle, showing the interactive graph, Story pane, and inspector at Birth in Québec](docs/images/chronicle-celine-dion-story.png)

## Research anything

People, families, places, institutions, events, inventions, species, objects, ideas,
and creative works. Study their origins, development, relationships, turning points,
historical context, and legacy or latest documented state.

## Rebuild the life. Play the story.

- **Automatic AI research:** enter a subject; Chronicle plans questions, searches sources, retrieves passages, and assembles chronological records.
- **Interactive histories:** explore events, relationships, works, and eras in a graph with date filters and adjustable detail.
- **Story playback:** read or listen to chapters while the timeline follows along; choose a perspective and reading speed.
- **Story mode:** follow a subject through time and place with narrated passages over a map that moves to each event's recorded locations.
- **Maps and media:** trace documented visits and connect images, music videos, and performances to their place in the story.
- **Academic exploration:** inspect source passages, compare competing claims, and preserve uncertain dates and coverage gaps.
- **Your research workspace:** annotate, save, import, and export projects; resize panes and switch themes.

| Follow places in Story mode | Explore music and videos |
| --- | --- |
| ![Céline Dion's Story mode at Birth in Québec, with a map of Charlemagne, chapter navigation, reading panel, and event cards](docs/images/chronicle-celine-dion-story-mode.png) | ![Céline Dion's updated video collection, including Dansons, with search and links to the original sources](docs/images/chronicle-celine-dion-videos.png) |

## Follow the story through time and place

Open **Story** in Explorer or choose **Story mode** above the Story pane. Read or
listen as the subject's timeline advances through passages and chapters. The map
follows the current event's recorded places, moving between city views and broader
regional context as the story unfolds.

This keeps the where alongside the when: follow a person's changing surroundings,
an institution's reach, or the geographic setting of a historical turning point.
In the Rome examples, the selected chapter and event cards move from the Republic
to the imperial system while the map keeps the city's geographic setting in view.

![Chronicle Story mode at From kings to the Republic, showing Rome on the map, the selected chapter, reading panel, and nearby events](docs/images/chronicle-rome-story-republic.png)

Pause, skip passages, seek along the story, choose a chapter, or adjust the voice
and speed. Switch off **Follow location** to explore the map yourself, or expand
the player to fullscreen. Opening Story mode keeps the current passage and playback.
The player includes a chapter sidebar, a reading panel, and nearby event cards.
Switch between Map, Timeline, and Graph while following the same selected passage;
open Places or Sources in the sidebar to inspect the story's geographic context
and evidence. Linked images appear with attribution when the project has them.

![Story mode at The Colosseum opens, with The imperial system chapter and the current event highlighted beside its narrative and map](docs/images/chronicle-rome-story-colosseum.png)

Research supplies geographic references from sources you can inspect under
**Map sources**. Several recorded locations are framed together; missing locations
use the wider story context. Regional markers are representative points, and the
modern base map does not reconstruct historical borders. [Story mode guide](docs/maps.md#story-mode).

## Start researching

Requires **Node.js 22+** and **Go 1.26+**.

```sh
npm install
npm run dev
```

Open [localhost:5173](http://localhost:5173). Connect your AI model and MCP search/retrieval
tools in **Settings**, then choose **New Research** and enter a subject.

[Setup](docs/setup.md) · [Research standards](docs/research-standards.md) · [Assistant/MCP tools](docs/assistant-tools.md) · [Maps](docs/maps.md) · [Speech](docs/speech.md) · [Image credits](docs/assets.md)
