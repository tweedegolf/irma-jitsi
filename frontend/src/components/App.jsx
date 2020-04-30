import React, { useState } from 'react';
import CssBaseline from '@material-ui/core/CssBaseline';
import Alert from '@material-ui/lab/Alert';
import CircularProgress from '@material-ui/core/CircularProgress';

import Jitsi from './Jitsi';

const JITSI_HOST = process.env.JITSI_HOST;
const BACKEND_HOST = process.env.BACKEND_HOST;

// For now we use a single room for the frontend.
// If you wish to add support for multiple rooms, add an appropriate user interface
// and submit a merge request.
const ROOM = process.env.ROOM;

const App = ({ startJitsiDisclosure }) => {
	const [error, setError] = useState(false);
	const [state, setState] = useState(null);

	return (
		<CssBaseline>
			{error && <Alert severity="error">Er ging iets mis, probeer het later nog eens.</Alert>}
			{state ? (
				<Jitsi
					domain={JITSI_HOST}
					jwt={state.jwt}
					roomName={state.room}
					userName={state.name}
					containerStyles={{ width: "auto", height: "auto" }} /* Use empty container styles */
					loadingComponent={<div className="content"><h1>IRMA + Jitsi</h1><CircularProgress /></div>}
				/>
			) : (
					<div className="content">
						<h1>IRMA + Jitsi</h1>
						<p>
							Je staat op het punt om deel te nemen aan een gesprek.
							Hiervoor dien je eerst jezelf te authenticeren.<br />
							Druk op de knop, scan de QR-code met de IRMA app, en je neemt direct deel aan het gesprek.
						</p>
						<button onClick={async () => {
							try {
								const { room, name, jwt } = await startJitsiDisclosure(`https://${BACKEND_HOST}`, ROOM);
								setState({ room, name, jwt });
							} catch (e) {
								if (e != 'CANCELLED') {
									console.error(e);
									setError(true);
								}
							}
						}}>
							Start authenticatie
					</button>
					</div>
				)}
		</CssBaseline>
	);
};

export default App;
