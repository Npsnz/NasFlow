# Implementation Plan - TaskFlow

TaskFlow is a production-grade personal and work task management web application designed to look, feel, and function like Asana or Linear. It features a local-first SQLite database using Go, Gin, GORM, and a modern Vue 3 interface using Pinia, Vue Router, and Tailwind CSS.

---

## User Review Required

We will implement the exact project layout, GORM models, REST endpoints, KanbanBoard, TaskDrawer, CalendarView, AllTasksView, Cmd+K command palette, recurring/overdue task services, BoardView stats, SettingsView, mobile responsiveness, Login/Registration/JWT settings, and detailed UX requirements as requested.

> [!IMPORTANT]
> - **GORM Models & Seeding**:
>   - GORM models will match all specified fields (including LexoRank-style `SortOrder` floats, `IsRecurring` booleans, and password hashes).
>   - Upon registering a new user, the backend will automatically seed 3 default workspaces: "งาน" (briefcase, #534AB7), "ส่วนตัว" (home, #1D9E75), "สุขภาพ" (heart, #D85A30).
> - **JWT Authentication & Rotation**:
>   - Short-lived Access Token (15m) + Long-lived Refresh Token (30d) stored in httpOnly cookie.
>   - Axios 401 interceptor rotates tokens and retries requests automatically.
> - **Detailed UX Interactions**:
>   - **Toasts**: Top-right placement, auto-dismiss 3s (Green: success, Red: error, Blue: info).
>   - **Loading States**: Skeleton loaders for task cards (no spinners for lists) and button spinners for form submissions. Optimistic updates on task movements.
>   - **Empty States**: Customized illustrations, helpful text, and CTA buttons. Search results display "ไม่พบงานที่ค้นหา" with filter clearing actions.
>   - **Keyboard Shortcuts**:
>     - `Cmd+K`/`Ctrl+K` → Quick Add Command Palette
>     - `Cmd+F`/`Ctrl+F` → Focus Search
>     - `E` (when task focused) → Edit focused task (opens Drawer)
>     - `D` (when task focused) → Toggle complete/done status
>     - `1` / `2` / `3` (when task focused) → Set priority to low, medium, high
>   - **Confirmation Modals**: Delete warning dialogs for Tasks ("ลบงานนี้หรือไม่? ไม่สามารถกู้คืนได้") and Workspaces containing active tasks.

---

## Proposed Changes

Here is the exact layout that we will create:

### Backend Architecture

```
backend/
├── main.go
├── go.mod
├── config/
│   └── config.go          # env vars, JWT secret, DB path
├── middleware/
│   ├── auth.go            # JWT validation middleware
│   └── cors.go
├── models/
│   ├── user.go
│   ├── workspace.go
│   ├── task.go
│   ├── tag.go
│   └── comment.go
├── handlers/
│   ├── auth.go
│   ├── workspace.go
│   ├── task.go
│   ├── tag.go
│   └── sse.go
├── services/
│   └── task_service.go    # business logic: recurring tasks, due date calc
└── database/
    └── db.go              # GORM setup + auto-migrate
```

- **[NEW] [models/user.go](file:///c:/Code/Workshops/NasFlow/backend/models/user.go)**: Matches the `User` struct exactly.
- **[NEW] [models/workspace.go](file:///c:/Code/Workshops/NasFlow/backend/models/workspace.go)**: Matches the `Workspace` struct.
- **[NEW] [models/task.go](file:///c:/Code/Workshops/NasFlow/backend/models/task.go)**: Matches `Task`.
- **[NEW] [models/tag.go](file:///c:/Code/Workshops/NasFlow/backend/models/tag.go)**: Matches `Tag`.
- **[NEW] [models/comment.go](file:///c:/Code/Workshops/NasFlow/backend/models/comment.go)**: Matches `Comment`.
- **[NEW] [database/db.go](file:///c:/Code/Workshops/NasFlow/backend/database/db.go)**: DB initialization, SQLite connection, auto-migration, and helper to seed default workspaces when a user is created.
- **[NEW] [config/config.go](file:///c:/Code/Workshops/NasFlow/backend/config/config.go)**: Port, SQLite DB path, JWT secret configurations.
- **[NEW] [middleware/auth.go](file:///c:/Code/Workshops/NasFlow/backend/middleware/auth.go)**: JWT middleware retrieving token from httpOnly cookie.
- **[NEW] [middleware/cors.go](file:///c:/Code/Workshops/NasFlow/backend/middleware/cors.go)**: CORS settings allowing credential forwarding.
- **[NEW] [handlers/auth.go](file:///c:/Code/Workshops/NasFlow/backend/handlers/auth.go)**: Handlers for `/api/auth/register`, `/api/auth/login`, `/api/auth/logout`, and `/api/auth/me`.
- **[NEW] [handlers/workspace.go](file:///c:/Code/Workshops/NasFlow/backend/handlers/workspace.go)**: CRUD for workspaces, including sorting order reordering.
- **[NEW] [handlers/task.go](file:///c:/Code/Workshops/NasFlow/backend/handlers/task.go)**: CRUD, filtering, column movement, reordering, subtasks, and completion checks.
- **[NEW] [handlers/tag.go](file:///c:/Code/Workshops/NasFlow/backend/handlers/tag.go)**: CRUD for tags.
- **[NEW] [handlers/sse.go](file:///c:/Code/Workshops/NasFlow/backend/handlers/sse.go)**: SSE subscription broker for real-time live boards.
- **[NEW] [services/task_service.go](file:///c:/Code/Workshops/NasFlow/backend/services/task_service.go)**: Recur task handling (generating next task after completion: daily, weekly, monthly, weekdays).
- **[MODIFY] [main.go](file:///c:/Code/Workshops/NasFlow/backend/main.go)**: Setup router and map endpoints.

---

### Frontend Architecture

```
frontend/
├── package.json
├── vite.config.ts
├── tailwind.config.js
├── src/
│   ├── main.ts
│   ├── App.vue
│   ├── router/index.ts
│   ├── stores/
│   │   ├── auth.ts
│   │   ├── tasks.ts
│   │   ├── workspace.ts
│   │   └── ui.ts          # sidebar state, active view
│   ├── api/
│   │   └── client.ts      # axios instance with interceptors
│   ├── components/
│   │   ├── layout/
│   │   │   ├── AppSidebar.vue
│   │   │   ├── AppTopbar.vue
│   │   │   └── AppLayout.vue
│   │   ├── board/
│   │   │   ├── KanbanBoard.vue
│   │   │   ├── KanbanColumn.vue
│   │   │   └── TaskCard.vue
│   │   ├── task/
│   │   │   ├── TaskDrawer.vue     # slide-in panel for task detail
│   │   │   ├── TaskForm.vue
│   │   │   ├── TaskList.vue
│   │   │   └── TaskFilters.vue
│   │   ├── calendar/
│   │   │   └── CalendarView.vue
│   │   └── ui/
│   │       ├── BaseModal.vue
│   │       ├── BaseBadge.vue
│   │       ├── BaseButton.vue
│   │       └── ToastNotification.vue
│   └── views/
│       ├── LoginView.vue
│       ├── BoardView.vue
│       ├── CalendarView.vue
│       ├── AllTasksView.vue
│       └── SettingsView.vue
```

---

## Verification Plan

### Automated tests:
```bash
cd backend
go test -v ./...
```

### Manual verification:
- Connect two browsers, verify immediate SSE synchronization on dragging cards.
- Verify automatic generation of next task instance for recurring tasks.
- Verify dashboard numbers updating correctly when tasks are completed or moved.
