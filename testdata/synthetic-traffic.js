import http from 'k6/http';
import { check, sleep } from 'k6';

const logins = ["admin@example.com", "admin2@example.com", "admin3@example.com"];

export default function () {
    const loginIndex = __ITER % logins.length;
    const url = 'http://localhost:8080/handle';
    const payload = JSON.stringify({
        "action": "auth",
        "auth": {
            "email": logins[loginIndex],
            "password": "verysecret"
        }
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    const res = http.post(url, payload, params);

    check(res, {
        'is status 202': (r) => r.status === 202,
    });

    sleep(0.3);
}