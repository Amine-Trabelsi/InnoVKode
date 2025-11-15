"""Collection of API routers grouped by feature area."""

from . import (
    admissions,
    ai_module,
    dashboard,
    dean,
    deadlines,
    dorms,
    events,
    hr,
    library,
    meta,
    rooms,
    schedule,
    support,
    teaching,
    users,
    visa,
    exams,
)

ROUTERS = [
    meta.router,
    users.router,
    schedule.router,
    exams.router,
    deadlines.router,
    rooms.router,
    events.router,
    ai_module.router,
    dean.router,
    admissions.router,
    teaching.router,
    hr.router,
    dashboard.router,
    dorms.router,
    library.router,
    support.router,
    visa.router,
]

