const API_ROOT =
    process.env.REACT_APP_BACKEND_URL ?? '/api';

function serialize(object) {
    const params = [];

    for (const param in object) {
        if (Object.hasOwnProperty.call(object, param) && object[param] != null) {
            params.push(`${param}=${encodeURIComponent(object[param])}`);
        }
    }

    return params.join('&');
}

let token = null;


const agent = async (url, body, method = 'GET', controller) => {
    const headers = new Headers();

    if (body) {
        headers.set('Content-Type', 'application/json');
    }

    if (token) {
        headers.set('Authorization', `Token ${token}`);
    }

    const response = await fetch(`${API_ROOT}${url}`, {
        method,
        headers,
        body: body ? JSON.stringify(body) : undefined,
        signal: controller,
    });
    let result;

    try {
        result = await response.json();
    } catch (error) {
        result = { errors: { [response.status]: [response.statusText] } };
    }

    if (!response.ok) throw result;

    return result;
};

const requests = {
    /**
     * Send a DELETE request
     *
     * @param {String} url The endpoint
     * @returns {Promise<Object>}
     */
    del: (url) => agent(url, undefined, 'DELETE'),
    /**
     * Send a GET request
     *
     * @param {String} url The endpoint
     * @param {Object} [query={}] URL parameters
     * @param {Number} [query.limit=10]
     * @param {Number} [query.page]
     * @param {String} [query.author]
     * @param {String} [query.tag]
     * @param {String} [query.favorited]
     * @returns {Promise<Object>}
     */
    get: (url, query = {}) => {
        if (Number.isSafeInteger(query?.page)) {
            query.limit = query.limit ? query.limit : 10;
            query.offset = query.page * query.limit;
        }
        delete query.page;
        const isEmptyQuery = query == null || Object.keys(query).length === 0;

        return agent(isEmptyQuery ? url : `${url}?${serialize(query)}`);
    },
    /**
     * Send a PUT request
     *
     * @param {String} url The endpoint
     * @param {Record<string, unknown>} body The request's body
     * @returns {Promise<Object>}
     */
    put: (url, body) => agent(url, body, 'PUT'),
    /**
     * Send a POST request
     *
     * @param {String} url The endpoint
     * @param {Record<string, unknown>} body The request's body
     * @returns {Promise<Object>}
     */
    post: (url, body, controller) => agent(url, body, 'POST', controller?.signal),
};

const Auth = {

    current: () => requests.get('/user'),

    login: (username, password) =>
        requests.post('/login', { username, password }),

    register: (username, email, password) =>
        requests.post('/users', { user: { username, email, password } }),

    save: (user) => requests.put('/user', { user }),
};

const Msg = {
    longpollMsg: (controller) => requests.post('/longpollMsg', {}, controller),

    sendMsg: (msg) => requests.post('sendMsg', msg)
};

const User = {
    getUserInfo: (username) => requests.get("userInfo", { username })
}

export default {
    Auth,
    Msg,
    User,
    setToken: (_token) => {
        token = _token;
    },
};