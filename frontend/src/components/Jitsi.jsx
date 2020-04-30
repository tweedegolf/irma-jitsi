import React, { createRef, useState, useEffect } from 'react';

// Heavily inspired by react-jutsu.

const Jitsi = ({
    roomName,
    userName,
    domain = 'meet.jit.si',
    password,
    jwt,
    loadingComponent,
    containerStyles,
    jitsiContainerStyles
}) => {
    const [loading, setLoading] = useState(true)
    const ref = createRef()

    const containerStyle = {
        width: '800px',
        height: '400px'
    }

    const jitsiContainerStyle = {
        display: loading ? 'none' : 'block',
        width: '100%',
        height: '100%'
    }

    const startConference = () => {
        try {
            const api = new JitsiMeetExternalAPI(domain, { roomName, parentNode: ref.current, jwt })
            api.addEventListener('videoConferenceJoined', () => {
                console.info(`${userName} has entered ${roomName}`)
                setLoading(false)
                api.executeCommand('displayName', userName)
                if (password) {
                    api.executeCommand('password', password)
                }
            })

            api.addEventListener('displayNameChange', () => {
                if (api.getDisplayName() !== userName) {
                    api.executeCommand('displayName', userName);
                }
            })
        } catch (error) {
            console.error('Failed to load Jitsi API', error)
        }
    }

    useEffect(() => {
        if (window.JitsiMeetExternalAPI) startConference()
        else console.error('Jitsi Meet API script not loaded')
    }, [])

    return (
        <div
            style={{ ...containerStyle, ...containerStyles }}
        >
            {loading && (loadingComponent || <p>Loading ...</p>)}
            <div
                id='jitsi-container'
                ref={ref}
                style={{ ...jitsiContainerStyle, ...jitsiContainerStyles }}
            />
        </div>
    )
};

export default Jitsi;