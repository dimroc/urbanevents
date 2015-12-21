import React, { Component } from 'react';
import { Link } from 'react-router';
import styles from './styles';

export default class Geoevent extends Component {
  render() {
    let geoevent = this.props.geoevent;
    var image = null
    console.log(geoevent.mediaUrl)
    if (geoevent.mediaType != "text") {
      image = <img src={geoevent.mediaUrl}/>
    }

    let className = styles.geoevent + " " + geoevent.mediaType;
    className += " uk-width-1-1 uk-width-medium-1-3 uk-panel uk-panel-box uk-panel-box-secondary"
    return <div className={className}>
      {geoevent.neighborhoods.map((hood) => {
        return <div className={styles.hood}>{hood}</div>
      })}
      <div className="uk-panel-title">{geoevent.fullName}</div>
      {image}
      <div>{geoevent.text}</div>
    </div>
  }
}

Geoevent.propTypes = {
  geoevent: React.PropTypes.object.isRequired
};

//<div className="full-name">{geoevent.fullName}</div>
//<div className="text">{geoevent.text}</div>
//<div className="media-url"><a href={geoevent.mediaUrl} target="_blank">{geoevent.mediaUrl}</a></div>
//<div className="neighborhoods">{geoevent.neighborhoods}</div>
//<div className="service">{geoevent.service}</div>
//<div className="created-at">{geoevent.createdAt}</div>
