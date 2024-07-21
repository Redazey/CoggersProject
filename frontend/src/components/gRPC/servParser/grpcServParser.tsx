import {servParser} from '../../../protos/servParser';
import { google } from '../../../protos/google/protobuf/empty';
import { getUUID } from '../../jwt/tokenUtils';

const clientServParser = new servParser.ServParserServiceClient('http://localhost:10000');
const uuid = getUUID().uuid;

export interface ServerInfo {
	Adress: string,
	Name: string,
	Version: string,
	MaxOnline: number,
	Online: number,
}

export function getServersInfo(): Promise<ServerInfo[]> {
    return new Promise((resolve, reject) => {
        clientServParser.GetServersInfo(new google.protobuf.Empty, {"uuid": uuid}, (err, response) => {
            if (err) {
                reject(err);
            } else if (response && response.serversInfo) {
                resolve(response.serversInfo);
            } else {
                reject(new Error('No servers info found'));
            }
        });
    });
}
