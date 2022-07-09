const baseUrl = new URL("/api", window.location);

export const isLoggedIn = () => {
    return document.cookie.includes("ddet_session");
};

export const logout = () => {
    document.cookie = `ddet_session=; Max-Age=0; Path=/; Domain=${document.location.hostname}; expires=Thu, 01 Jan 1970 00:00:00 GMT`;
    document.location.reload();
};

export const get = async (url) => {
    return request(url);
};
export const post = async (url, body) => {
    return request(url, body);
};

export const request = async (url, body, options = {}) => {
    const headers = {
        accepts: "application/json",
        ...options.headers,
    };

    if (!options.method && body) {
        options.method = "POST";
    }

    if (options.method && options.method !== "GET") {
        headers["content-type"] = "application/json";

        if (body !== undefined && typeof body !== "string") {
            options.body = JSON.stringify(body);
        }
    }

    return fetch(new URL(url, baseUrl), {
        ...options,
        headers,
        body,
    }).then(
        async (res) => {
            if (!res.ok) {
                if (res.status === 401) {
                    logout();
                    return {};
                }
                return res.json().then(({ error }) => {
                    console.error(e);
                    console.error("api error: " + error);
                    throw new Error(error);
                });
            }
            if (res.status == 201) {
                return res;
            }
            return res.json();
        },
        (e) => {
            console.error(e);
            console.error("network unreachable: " + e.message);
            throw e;
        },
    );
};
