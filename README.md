# NasFlow — Task Management System

A production-grade personal task management web application inspired by Linear and Asana. Built with a Go backend and Vue 3 frontend, designed for local self-hosting with a modern, minimal UI.

---

## Features

**Task Management**
- Kanban board with drag-and-drop task reordering across status columns
- Full task lifecycle: Todo → In Progress → Done
- Priority levels: Low, Medium, High, Urgent with visual indicators
- Due dates, recurring tasks (daily / weekly / monthly / weekdays)
- Subtasks, comments, and tag system with custom colors

**Workspaces**
- Multiple isolated workspaces (Work, Personal, Health — seeded on registration)
- Per-workspace task filtering and board views

**Real-time**
- Server-Sent Events (SSE) for live board updates across browser tabs
- Overdue task alerts pushed automatically

**UX**
- `Ctrl+K` — Quick Add command palette with smart syntax parsing
- `Ctrl+F` — Focus search bar
- Keyboard shortcuts on focused tasks: `E` (edit), `D` (toggle done), `1/2/3` (priority)
- Dark / Light / System theme
- Thai and English language support
- Responsive design (mobile, tablet, desktop)

**Authentication**
- JWT with token rotation: 15-minute access token + 30-day refresh token via httpOnly cookie
- Automatic token refresh via Axios interceptor

---

## Tech Stack

| Layer | Technology |
|---|---|
| Backend | Go 1.25, Gin, GORM, SQLite |
| Frontend | Vue 3, TypeScript, Pinia, Vue Router |
| Styling | Tailwind CSS v3 |
| Auth | JWT (golang-jwt/jwt v5) |
| Real-time | Server-Sent Events |
| Icons | Lucide Vue Next |

---

## Project Structure

```
NasFlow/
├── backend/
│   ├── main.go
│   ├── config/         # Environment variables and JWT config
│   ├── database/       # GORM setup, SQLite connection, auto-migrate
│   ├── handlers/       # HTTP route handlers (auth, tasks, workspaces, tags, SSE)
│   ├── middleware/     # JWT auth middleware, CORS
│   ├── models/         # GORM models (User, Task, Workspace, Tag, Comment)
│   └── services/       # Business logic (recurring tasks, overdue detection)
└── frontend/
    └── src/
        ├── api/            # Axios instance with interceptors
        ├── components/
        │   ├── board/      # KanbanBoard, KanbanColumn, TaskCard
        │   ├── calendar/   # CalendarView
        │   ├── layout/     # AppSidebar, AppTopbar, AppLayout
        │   ├── task/       # TaskDrawer, TaskForm, TaskList
        │   └── ui/         # BaseButton, BaseBadge, BaseModal, Toast
        ├── router/         # Vue Router with auth guards
        ├── stores/         # Pinia stores (auth, tasks, workspace, ui)
        └── views/          # BoardView, CalendarView, AllTasksView, SettingsView, LoginView
```

---

## Getting Started

### Prerequisites

- [Go 1.21+](https://go.dev/dl/)
- [Node.js 18+](https://nodejs.org/)
- A C compiler (for SQLite — e.g. `gcc` on Linux/Mac, [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) on Windows)

### 1. Clone

```bash
git clone https://github.com/your-username/nasflow.git
cd nasflow
```

### 2. Start the Backend

```bash
cd backend
go run main.go
```

The API server starts at **http://localhost:8080**. The SQLite database (`taskflow.db`) is created automatically on first run.

### 3. Start the Frontend

```bash
cd frontend
npm install
npm run dev
```

The frontend dev server starts at **http://localhost:5173**.

### 4. Open the App

Navigate to **http://localhost:5173** and register a new account. Three default workspaces (Work, Personal, Health) are created automatically.

---

## API Endpoints

### Auth
| Method | Path | Description |
|---|---|---|
| `POST` | `/api/auth/register` | Register new account |
| `POST` | `/api/auth/login` | Login |
| `POST` | `/api/auth/logout` | Logout (clears cookie) |
| `GET` | `/api/auth/me` | Get current user |
| `POST` | `/api/auth/refresh` | Refresh access token |
| `PUT` | `/api/auth/profile` | Update profile |
| `DELETE` | `/api/auth/delete` | Delete account |

### Tasks
| Method | Path | Description |
|---|---|---|
| `GET` | `/api/tasks` | List tasks (filterable) |
| `POST` | `/api/tasks` | Create task |
| `GET` | `/api/tasks/:id` | Get task details |
| `PUT` | `/api/tasks/:id` | Update task |
| `DELETE` | `/api/tasks/:id` | Delete task |
| `PUT` | `/api/tasks/:id/status` | Update task status |
| `POST` | `/api/tasks/:id/complete` | Mark complete (triggers recurring) |
| `PUT` | `/api/tasks/reorder` | Reorder tasks (LexoRank-style) |
| `GET` | `/api/tasks/overdue` | Get overdue tasks |

### Workspaces, Tags, Comments
| Method | Path | Description |
|---|---|---|
| `GET/POST` | `/api/workspaces` | List / Create |
| `PUT/DELETE` | `/api/workspaces/:id` | Update / Delete |
| `GET/POST` | `/api/tags` | List / Create |
| `PUT/DELETE` | `/api/tags/:id` | Update / Delete |
| `POST` | `/api/tasks/:id/comments` | Add comment |
| `DELETE` | `/api/comments/:id` | Delete comment |
| `GET` | `/api/stats` | Dashboard statistics |
| `GET` | `/api/sse` | SSE stream for real-time updates |

---

## Keyboard Shortcuts

| Shortcut | Action |
|---|---|
| `Ctrl+K` | Open Quick Add command palette |
| `Ctrl+F` | Focus search bar |
| `E` | Open Task Drawer (requires focused task) |
| `D` | Toggle task done/todo (requires focused task) |
| `1` | Set priority: Low (requires focused task) |
| `2` | Set priority: Medium (requires focused task) |
| `3` | Set priority: High (requires focused task) |
| `Esc` | Close drawer / modal |

To focus a task, click on its card. A blue ring indicates the focused task.

---

## Quick Add Syntax

The `Ctrl+K` command palette supports inline smart parsing:

```
Buy groceries #personal !high @tomorrow every week
```

| Syntax | Example | Result |
|---|---|---|
| `#workspace` | `#งาน` | Assign to workspace |
| `!priority` | `!high`, `!ด่วน` | Set priority |
| `@date` | `@today`, `@วันนี้`, `@monday` | Set due date |
| `every ...` | `every day`, `ทุกสัปดาห์` | Set recurrence |

---

## Environment Configuration

The backend reads from environment variables with sensible defaults:

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8080` | HTTP server port |
| `DB_PATH` | `./taskflow.db` | SQLite database file path |
| `JWT_SECRET` | *(generated)* | Secret key for JWT signing |
| `JWT_ACCESS_TTL` | `15m` | Access token lifetime |
| `JWT_REFRESH_TTL` | `720h` | Refresh token lifetime (30 days) |

---

## Production Build

```bash
# Build frontend
cd frontend
npm run build
# Output: frontend/dist/

# Build backend binary
cd backend
go build -o nasflow main.go
./nasflow
```

For production, serve the `frontend/dist/` directory via a reverse proxy (nginx, Caddy) pointing to the Go server at port 8080.

---

## License

MIT
