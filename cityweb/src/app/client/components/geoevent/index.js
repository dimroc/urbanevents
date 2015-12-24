import React, { Component } from 'react';
import { Link } from 'react-router';
import styles from './styles';
import moment from 'moment';

export default class Geoevent extends Component {
  render() {
    let geoevent = this.props.geoevent;
    let image = null
    let timeAgo = moment(geoevent.createdAt).fromNow();

    if (geoevent.mediaType != "text") {
      image = <img src={geoevent.mediaUrl}/>
    }

    let className = styles.geoevent + " " + geoevent.mediaType;
    className += " uk-width-1-1"
    return <a className={className} href={geoevent.link} target="_blank">
      <div className={styles.hood}>
        {geoevent.neighborhoods.map((hood) => {
          return <span key={ hood }>{hood}</span>
        })}
      </div>
      {image}
      <div className={styles.meta}>
        <h3>{geoevent.fullName}</h3>
        <label className="uk-badge">{timeAgo}</label>
        <p>{geoevent.text}</p>
      </div>
    </a>
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
