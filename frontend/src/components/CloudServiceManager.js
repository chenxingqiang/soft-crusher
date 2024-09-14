// File: frontend/src/components/CloudServiceManager.js

import React, { useState, useEffect } from 'react'
import axios from 'axios'

function CloudServiceManager() {
    const [cloudProvider, setCloudProvider] = useState('')
    const [clusterName, setClusterName] = useState('')
    const [nodeCount, setNodeCount] = useState(3)
    const [status, setStatus] = useState('')
    const [progress, setProgress] = useState(0)

    const handleDeploy = async () => {
        try {
            setStatus('Deploying...')
            setProgress(0)

            const response = await axios.post(
                '/api/deploy-to-cloud',
                {
                    cloudProvider,
                    clusterName,
                    nodeCount
                },
                {
                    headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
                }
            )

            setStatus('Deployment initiated. Checking progress...')
            checkDeploymentProgress()
        } catch (error) {
            setStatus('Deployment failed: ' + error.response.data)
        }
    }

    const checkDeploymentProgress = async () => {
        try {
            const response = await axios.get('/api/deployment-progress', {
                headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
            })

            setProgress(response.data.progress)
            setStatus(response.data.status)

            if (response.data.progress < 100) {
                setTimeout(checkDeploymentProgress, 5000)
            }
        } catch (error) {
            setStatus('Error checking deployment progress: ' + error.response.data)
        }
    }

    return (
        <div>
            <h2>Cloud Service Manager</h2>
            <select value={cloudProvider} onChange={(e) => setCloudProvider(e.target.value)}>
                <option value="">Select Cloud Provider</option>
                <option value="aliyun">Alibaba Cloud</option>
                <option value="aws">Amazon Web Services</option>
            </select>
            <input
                type="text"
                placeholder="Cluster Name"
                value={clusterName}
                onChange={(e) => setClusterName(e.target.value)}
            />
            <input
                type="number"
                placeholder="Node Count"
                value={nodeCount}
                onChange={(e) => setNodeCount(parseInt(e.target.value))}
            />
            <button onClick={handleDeploy}>Deploy to Cloud</button>
            <p>{status}</p>
            <progress value={progress} max="100"></progress>
        </div>
    )
}

export default CloudServiceManager
