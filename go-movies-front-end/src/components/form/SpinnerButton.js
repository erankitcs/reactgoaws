import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faCircleNotch } from '@fortawesome/free-solid-svg-icons';


const SpinnerButton = (props) => {
    console.log(props.state);
    const className = props.state ? props.className + ' disabled' : props.className;
    return (
      <button
        type={props.type}
        className={className}
        onClick={props.onClick}
      >
        {props.state ? (
          <>
            <FontAwesomeIcon icon={faCircleNotch} spin />
          </>
        ) : (
          <>
          <FontAwesomeIcon icon={props.faIcon} />
          </>
        )}
      </button>
    );
  };

export default SpinnerButton;