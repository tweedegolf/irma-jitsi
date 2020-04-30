import axios from 'axios';
import { handleSession } from '@privacybydesign/irmajs';
import React from 'react';
import ReactDOM from 'react-dom';

import App from './components/App';

const veiligJitsi = {
    start: async ({ url, room }) => {
        const irmaResponse = await axios.get(`${url}/session`, { params: { room } });

        if (irmaResponse.status !== 200) {
            throw new Error(`Failed to create session with status ${irmaResponse.status}`);
        }

        const { sessionPtr, trustedFacts } = irmaResponse.data;

        await handleSession(sessionPtr, { language: 'nl' });

        const discloseResponse = await axios.get(`${url}/disclose`, { params: { trustedFacts } });

        if (discloseResponse.status !== 200) {
            throw new Error(`Failed to fetch disclosure with status ${discloseResponse.status}`);
        }

        return discloseResponse.data; // Room, name, jwt
    },
};

window.addEventListener('load', () => {
    const container = window.document.getElementById('container');

    ReactDOM.render(
        <App startJitsiDisclosure={(url, room) => veiligJitsi.start({ url, room })} />,
        container
    );
});
