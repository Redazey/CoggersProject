import {auth} from '../../../protos/auth';
import { getUUID } from '../../jwt/tokenUtils';

const clientAuth = new auth.AuthServiceClient('http://localhost:8080')
const uuid = getUUID();

export function Registration(
    name: string,
    email: string,
    password: string,
    roleId: number,
    birthdate: string,
    photourl: string,
    push: boolean
): Promise<string> {
    return new Promise((resolve, reject) => {
        const request = new auth.RegistrationRequest();
        request.name = name
        request.email = email
        request.password = password
        request.roleId = roleId
        request.birthdate = birthdate
        request.photourl = photourl
        request.push = push

        clientAuth.Registration(request, {"uuid": uuid}, (err, response) => {
            if (err) {
                reject(err);
            } else if (response && response.key) {
                resolve(response.key);
            } else {
                reject(new Error('No comments found'));
            }
        });
    });
}

export function Login(email: string, password: string): Promise<string> {
    return new Promise((resolve, reject) => {
        const request = new auth.RegistrationRequest();
        request.email = email
        request.password = password
        
        clientAuth.Registration(request, {"uuid": uuid}, (err, response) => {
            if (err) {
                reject(err);
            } else if (response && response.key) {
                resolve(response.key);
            } else {
                reject(new Error('No comments found'));
            }
        });
    });
}
