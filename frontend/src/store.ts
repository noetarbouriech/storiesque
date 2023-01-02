import { writable, type Writable } from 'svelte/store';
import { browser } from '$app/environment';

export type User = {
    id: number,
    username: string,
    email: string,
    is_admin: boolean,
    has_img: boolean
};

let userStore = writable<User>(JSON.parse("{\"username\":\"\",\"email\":\"\"}"));
if (browser) {
    const user = localStorage.getItem("user");
    if (user != null) userStore = writable<User>(JSON.parse(user));
    userStore.subscribe(value => {
        localStorage.setItem("user", JSON.stringify(value));
    })
}
export { userStore };