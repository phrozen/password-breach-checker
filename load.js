import http from 'k6/http';

const randSHA1 = (size=40) => [...Array(size)].map(() => {
  return Math.floor(Math.random() * 16).toString(16)
}).join('');

export default function () {
    const url = 'http://localhost:3000/check';
    const payload = `hash=${randSHA1()}`;
    const params = {
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
    };
    http.post(url, payload, params);
}

// k6 run --vus 30 --duration 60s load.js