import { useAppContext } from '../App.jsx'

function VolumeControl () {
  const { volume, handleVolumeChange, isSoundEnabled, toggleSound } =
    useAppContext()

  // スライダーが変更されたときに呼ばれる関数
  const handleChange = event => {
    handleVolumeChange(parseInt(event.target.value, 10))
  }

  return (
    <div className='volume-control-wrapper'>
      <div className='volume-control'>
        <label htmlFor='volumeSlider'>音量: </label>
        <input
          type='range'
          id='volumeSlider'
          min='0'
          max='100'
          value={volume}
          onChange={handleChange}
          aria-label='音量調節スライダー'
          disabled={!isSoundEnabled}
        />
        <span id='volumeValue' className='volume-display'>
          {volume}%
        </span>

        <input
          type='checkbox'
          id='soundToggleCheckbox'
          checked={isSoundEnabled}
          onChange={toggleSound}
        />
        <label htmlFor='soundToggleCheckbox'>音を有効にする</label>
      </div>
    </div>
  )
}
export default VolumeControl
