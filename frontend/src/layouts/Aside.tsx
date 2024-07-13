import React, {useState, useEffect, useCallback} from "react";
import { Link } from 'react-router-dom';
import {ServerInfo, getServersInfo} from '../components/gRPC/servParser/grpcServParser'
import { google } from '../protos/google/protobuf/empty';

interface ServerCardProps {
    server: ServerInfo;
}

const ServerCard: React.FC<ServerCardProps> = ({ server }) => {
    return (
        <div className="servers__card">
            <Link to={`/server/${server.Id}`}><h2>{server.Name}</h2></Link>
            <table>
                <tr>
                    <td>онлайн:</td>
                    <td>{server.Online}/{server.MaxOnline}</td>
                </tr>
                <tr>
                    <td>адрес:</td>
                    <td>{server.Adress}</td>
                </tr>
                <tr>
                    <td>версия:</td>
                    <td>{server.Version}</td>
                </tr>
            </table>
        </div>
    );
}

const Aside: React.FC = () => {
    const [serversInfo, setServerInfo] = useState<Map<string, ServerInfo> | null>(null);

    const fetchServersInfo = useCallback(async () => {
        try {
            const ServParserResponse = await getServersInfo();
            setServerInfo(ServParserResponse);
        } catch (error) {
            console.error('Error fetching movie:', error);
        }
    }, [new google.protobuf.Empty]);

    useEffect(() => {
        fetchServersInfo();
    }, [fetchServersInfo]);

    if (!serversInfo) {
        return <div>Загрузка...</div>
    }

    return (
        <>
            <div className="servers">
                {Array.from(serversInfo).map(([key, value]) => (
                    <ServerCard key={key} server={value}/>
                ))};
            </div>
        </>
    );
};

export default Aside;
