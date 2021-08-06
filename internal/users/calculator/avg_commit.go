package calculator

import (
	"log"
	"signum-explorer-bot/internal/api/signum_api"
	"signum-explorer-bot/internal/config"
	"signum-explorer-bot/internal/database/models"
	"sync"
	"time"
)

func (c *Calculator) GetLastMiningIngo() signum_api.MiningInfo {
	c.RLock()
	lastMiningInfo := c.lastMiningInfo
	c.RUnlock()
	return lastMiningInfo
}

func (c *Calculator) readAverageCommitmentFromDB() {
	var avgCommitments []models.AverageCommitment
	result := c.db.Find(&avgCommitments)
	if result.Error != nil || result.RowsAffected == 0 {
		log.Printf("Error getting Average Commitments from DB: %v", result.Error)
		return
	}

	var sum float64
	for _, v := range avgCommitments {
		sum += v.AverageCommitment
	}
	c.lastMiningInfo.AverageCommitmentNQT = sum / float64(len(avgCommitments))
	log.Printf("Have loaded Average Commitment from DB: %v", c.lastMiningInfo.AverageCommitmentNQT)
}

func (c *Calculator) saveNewCommitment(newCommitment float64) {
	dbCommit := &models.AverageCommitment{
		AverageCommitment: newCommitment,
	}
	c.db.Save(dbCommit)
	log.Printf("Have got and saved new Average Commitment %v", newCommitment)

	// delete irrelevant data
	if config.SIGNUM_API.AVERAGING_QUANTITY < dbCommit.ID {
		c.db.Unscoped().Delete(models.AverageCommitment{}, "id <= ?", dbCommit.ID-config.SIGNUM_API.AVERAGING_QUANTITY)
	}
}

func (c *Calculator) StartAverageCommitmentListener(wg *sync.WaitGroup, shutdownChannel chan interface{}) {
	defer wg.Done()

	log.Printf("Start Average Commitment Listener")
	ticker := time.NewTicker(config.SIGNUM_API.GET_AVG_COMMIT_TIME)

	c.updateMiningInfo()
	for {
		select {
		case <-shutdownChannel:
			log.Printf("Average Commitment Listener received shutdown signal")
			ticker.Stop()
			return

		case <-ticker.C:
			c.updateMiningInfo()
		}
	}
}

func (c *Calculator) updateMiningInfo() {
	miningInfo, err := c.signumClient.GetMiningInfo()
	if err != nil {
		log.Printf("Error getting mining info: %v", err)
		return
	}
	newCommitment := miningInfo.AverageCommitmentNQT / 1e8
	c.saveNewCommitment(newCommitment)

	var count int64
	c.db.Model(&models.AverageCommitment{}).Count(&count)
	if count > 0 {
		c.Lock() // update global value
		prevCommitment := c.lastMiningInfo.AverageCommitmentNQT
		c.lastMiningInfo = *miningInfo
		c.lastMiningInfo.AverageCommitmentNQT = (prevCommitment*float64(count-1) + newCommitment) / float64(count)
		c.Unlock()
	}
}
