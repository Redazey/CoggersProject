import Cookies from 'js-cookie';

export function getUUID(): { uuid: string } {
    const userUUID = Cookies.get("userUUID");
    const uuid = userUUID ? userUUID : "";

    return { uuid };
}

export function getToken(): { token: string } {
    const jwt = Cookies.get("token");
    const token = jwt ? jwt : "";

    return { token };
}

