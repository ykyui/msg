import { openDB } from 'idb';

const dbPromise = openDB("storage", 1, {
    upgrade(db, oldVersion, newVersion, transaction) {
        if (!db.objectStoreNames.contains("message")) {
            db.createObjectStore('message', { keyPath: 'id' });
        }
    },
    blocked() {
        // …
    },
    blocking() {
        // …
    },
    terminated() {
        // …
    },
});

export async function get(key) {
    return (await dbPromise).get('message', key);
};
export async function set(key, val) {
    return (await dbPromise).put('message', val, key);
};
export async function del(key) {
    return (await dbPromise).delete('message', key);
};
export async function clear() {
    return (await dbPromise).clear('message');
};
export async function keys() {
    return (await dbPromise).getAllKeys('message');
};

export async function msgStorage() {
    return (await dbPromise).objectStoreNames("message")
}
export async function msgContent() {
    const s = (await msgStorage())
    const from = s.index("from").getAllKeys()
    const to = s.index("to").getAllKeys()
    return [...from, ...to]
}