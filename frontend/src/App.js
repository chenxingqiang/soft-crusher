import React from 'react'
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom'
import Header from './components/Header'
import Footer from './components/Footer'
import Home from './pages/Home'
import Dashboard from './pages/Dashboard'

const App = () => {
    return (
        <Router>
            <div className="app">
                <Header />
                <main>
                    <Switch>
                        <Route exact path="/" component={Home} />
                        <Route path="/dashboard" component={Dashboard} />
                    </Switch>
                </main>
                <Footer />
            </div>
        </Router>
    )
}

export default App
